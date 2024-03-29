package library

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/BrenekH/encodarr/controller"
	"github.com/google/uuid"
)

// NewManager return a new Manager.
func NewManager(logger controller.Logger, ds controller.LibraryManagerDataStorer, metadataReader MetadataReader, commandDecider CommandDecider) Manager {
	return Manager{
		logger:         logger,
		ds:             ds,
		metadataReader: metadataReader,
		commandDecider: commandDecider,
		videoFileser:   defaultVideoFileser{},
		fileRemover:    defaultFileRemover{},
		fileMover:      defaultFileMover{},
		fileStater:     defaultFileStater{},

		lastCheckedTimes:   make(map[int]time.Time),
		workerCompletedMap: make(map[int]bool),
	}
}

// Manager satisfies the conroller.LibraryManager interface.
type Manager struct {
	logger         controller.Logger
	ds             controller.LibraryManagerDataStorer
	metadataReader MetadataReader
	commandDecider CommandDecider
	videoFileser   videoFileser
	fileRemover    fileRemover
	fileMover      fileMover
	fileStater     fileStater

	// lastCheckedTimes is a map of Library ids and the last time that they were checked.
	lastCheckedTimes map[int]time.Time

	// workerCompletedMap is a map of Library ids and a boolean to indicate whether the goroutine that was spawned is finished
	workerCompletedMap map[int]bool
}

// Start starts the library manager without blocking the thread.
func (m *Manager) Start(ctx *context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if controller.IsContextFinished(ctx) {
				return
			}

			// Check all Libraries for required scans
			allLibraries, err := m.ds.Libraries()
			if err != nil {
				m.logger.Error("%v", err)
				time.Sleep(time.Second)
				continue
			}

			for _, lib := range allLibraries {
				t, ok := m.lastCheckedTimes[lib.ID]
				if !ok {
					m.lastCheckedTimes[lib.ID] = time.Unix(0, 0)
					t = m.lastCheckedTimes[lib.ID]
				}

				previousWorkerFinished, ok := m.workerCompletedMap[lib.ID]
				if !ok {
					m.workerCompletedMap[lib.ID] = true
					previousWorkerFinished = m.workerCompletedMap[lib.ID]
				}

				if time.Since(t) > lib.FsCheckInterval && previousWorkerFinished {
					m.logger.Debug("Initiating library (ID: %v) update", lib.ID)
					m.lastCheckedTimes[lib.ID] = time.Now()
					m.workerCompletedMap[lib.ID] = false

					wg.Add(1)
					go m.updateLibraryQueue(ctx, wg, lib)
				}
			}
			time.Sleep(time.Second)
		}
	}()
}

func (m *Manager) updateLibraryQueue(ctx *context.Context, wg *sync.WaitGroup, lib controller.Library) {
	defer wg.Done()
	defer func() { m.workerCompletedMap[lib.ID] = true }()

	// Locate video files
	discoveredVideos, err := m.videoFileser.VideoFiles(lib.Folder)
	if err != nil {
		m.logger.Error(err.Error())
		return
	}

	for _, videoFilepath := range discoveredVideos {
		// Respect context while iterating over discoveredVideos
		if controller.IsContextFinished(ctx) {
			break
		}

		// Check path against Library path masks
		maskedOut := false
		for _, v := range lib.PathMasks {
			if v == "" {
				m.logger.Trace("Skipping an empty path mask string")
				continue
			}
			if strings.Contains(videoFilepath, v) {
				m.logger.Debug("%v skipped because of a mask (%v)", videoFilepath, v)
				maskedOut = true
				break
			}
		}
		// Use the maskedOut variable to continue the iteration over discovered media and not the path masks
		if maskedOut {
			continue
		}

		pathDispatched, err := m.ds.IsPathDispatched(videoFilepath)
		if err != nil {
			m.logger.Error(err.Error())
			continue
		}

		if pathDispatched || lib.Queue.InQueuePath(controller.Job{Path: videoFilepath}) {
			continue
		}

		// Read file metadata from a MetadataReader
		fMetadata, err := m.metadataReader.Read(videoFilepath)
		if err != nil {
			m.logger.Error("Skipping %v because of error: %v", videoFilepath, err)
			continue
		}

		// Run a CommandDecider against the metadata to determine what FFMpeg command to run
		commandSlice, err := m.commandDecider.Decide(fMetadata, lib.CommandDeciderSettings)
		if err != nil {
			m.logger.Debug("Skipping %v because CommandDecider returned error: %v", videoFilepath, err)
			continue
		}

		// Save to Library queue
		job := controller.Job{
			UUID:     controller.UUID(uuid.NewString()),
			Path:     videoFilepath,
			Command:  commandSlice,
			Metadata: fMetadata,
		}
		lib.Queue.Push(job)
		m.logger.Info("Added %v to Library %v's queue", videoFilepath, lib.ID)

		m.ds.SaveLibrary(lib)
	}
}

// ImportCompletedJobs takes a list of completed jobs and imports them and their files into the system.
func (m *Manager) ImportCompletedJobs(jobs []controller.CompletedJob) {
	for _, cJob := range jobs {
		// Pop job from dispatched_jobs
		dJob, err := m.ds.PopDispatchedJob(cJob.UUID)
		if err != nil {
			m.logger.Error(err.Error())
			continue
		}

		// If job failed, log it, save the history entry to the history table, and continue iterating.
		if cJob.Failed {
			m.logger.Warn("Job for file %v failed: %v, %v", dJob.Job.Path, cJob.History.Warnings, cJob.History.Errors)
			if err = m.ds.PushHistory(cJob.History); err != nil {
				m.logger.Error(err.Error())
			}
			continue
		}

		//? Somewhere in here should be an evaluation from CommandDecider to detect if any plugins want to make more changes. If they do then the file should be placed in a cache location and not the og file location.

		filename := dJob.Job.Path

		// Remove old file
		if err = m.fileRemover.Remove(dJob.Job.Path); err != nil {
			failMessage := fmt.Sprintf("Failed to remove file '%v' because of error: %v", dJob.Job.Path, err)
			m.logger.Error(failMessage)

			// Set filename to a string with an extra encodarr extension
			fnExt := filepath.Ext(filename)
			i := strings.LastIndex(filename, fnExt)
			fnWoExt := filename[:i] + strings.Replace(filename[i:], fnExt, "", 1)
			filename = fmt.Sprintf("%v.encodarr%v", fnWoExt, fnExt)

			cJob.History.Warnings = append(cJob.History.Warnings, failMessage)
		}

		// Change filename to have file extension of InFile
		inFileExt := filepath.Ext(cJob.InFile)
		fnExt := filepath.Ext(filename)
		i := strings.LastIndex(filename, fnExt)
		fnWoExt := filename[:i] + strings.Replace(filename[i:], fnExt, "", 1)
		filename = fnWoExt + inFileExt

		// Move new file to old file location
		if err = m.fileMover.Move(cJob.InFile, filename); err != nil {
			failMessage := fmt.Sprintf("Failed to move file '%v' because of error: %v", dJob.Job.Path, err)
			m.logger.Error(failMessage)

			cJob.History.Errors = append(cJob.History.Errors, failMessage)
		}

		// Save history entry to histroy table
		if err = m.ds.PushHistory(cJob.History); err != nil {
			m.logger.Error(err.Error())
		}
	}
}

// LibrarySettings returns the current settings of each library in the data store.
func (m *Manager) LibrarySettings() ([]controller.Library, error) {
	libs, err := m.ds.Libraries()

	if err != nil {
		m.logger.Error(err.Error())
	}

	return libs, err
}

// PopNewJob returns and deletes a job from the library queues in order of priority.
func (m *Manager) PopNewJob() (controller.Job, error) {
	// Get every library from DataStorer (m.ds.Libraries())
	libs, err := m.ds.Libraries()
	if err != nil {
		m.logger.Error(err.Error())
		return controller.Job{}, err
	}

	// Sort libraries by decreasing order so that the libraries with the higher priority number dispatch jobs first.
	sort.Slice(libs, func(i, j int) bool {
		return libs[i].Priority > libs[j].Priority
	})

	// Loop through sorted slice looking for a job to return
	for _, l := range libs {
		for !l.Queue.Empty() {
			job, err := l.Queue.Pop()
			if err != nil {
				if err != controller.ErrEmptyQueue { // Forgoes logging about an empty queue
					m.logger.Debug("error while searching for job: %v", err)
				}
				continue
			}

			// Skip queue entry if there is an error while stating the file
			_, err = m.fileStater.Stat(job.Path)
			if err != nil {
				m.logger.Debug("skipping queue entry for %v because of error: %v", job.Path, err)
				continue
			}

			// Update library in datastore
			err = m.ds.SaveLibrary(l)
			if err != nil {
				m.logger.Error(err.Error())
			}

			return job, nil
		}
	}

	return controller.Job{}, fmt.Errorf("no available jobs")
}

// UpdateLibrarySettings loops through each entry in the provided map and applies the new settings
// if the key matches a valid library. However, it will not update the ID and Queue fields.
// If the key doesn't match a valid library, a brand new one with the provided settings is created.
func (m *Manager) UpdateLibrarySettings(libSettings map[int]controller.Library) {
	for k, v := range libSettings {
		lib, err := m.ds.Library(k)
		if err != nil {
			// Save brand new library with key as ID and value as library object
			v.ID = k
			v.Queue = controller.LibraryQueue{}
			v.CommandDeciderSettings = m.commandDecider.DefaultSettings()

			if err = m.ds.SaveLibrary(v); err != nil {
				m.logger.Error(err.Error())
			}
			continue
		}

		lib.Folder = v.Folder
		lib.Priority = v.Priority
		lib.FsCheckInterval = v.FsCheckInterval
		lib.PathMasks = v.PathMasks
		lib.CommandDeciderSettings = v.CommandDeciderSettings

		if err = m.ds.SaveLibrary(lib); err != nil {
			m.logger.Error(err.Error())
		}
	}
}

type defaultVideoFileser struct{}

func (d defaultVideoFileser) VideoFiles(dir string) ([]string, error) {
	return GetVideoFilesFromDir(dir)
}

type defaultFileRemover struct{}

func (d defaultFileRemover) Remove(path string) error {
	return os.Remove(path)
}

type defaultFileMover struct{}

func (d defaultFileMover) Move(from, to string) error {
	inputFile, err := os.Open(from)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %s", err)
	}

	outputFile, err := os.Create(to)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("couldn't open dest file: %s", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("writing to output file failed: %s", err)
	}

	// The copy was successful, so now delete the original file
	err = os.Remove(from)
	if err != nil {
		return fmt.Errorf("failed removing original file: %s", err)
	}

	return nil
}

type defaultFileStater struct{}

func (d defaultFileStater) Stat(path string) (fs.FileInfo, error) {
	return os.Stat(path)
}
