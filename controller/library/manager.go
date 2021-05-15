package library

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/BrenekH/encodarr/controller"
	"github.com/google/uuid"
)

func NewManager(logger controller.Logger, ds controller.LibraryManagerDataStorer) Manager {
	// TODO: Take a MetadataReader and a CommandDecider as arguments
	return Manager{
		logger: logger,
		ds:     ds,

		lastCheckedTimes:   make(map[int]time.Time),
		workerCompletedMap: make(map[int]bool),
	}
}

type Manager struct {
	logger controller.Logger
	ds     controller.LibraryManagerDataStorer

	metadataReader MetadataReader
	commandDecider CommandDecider

	// lastCheckedTimes is a map of Library ids and the last time that they were checked.
	lastCheckedTimes map[int]time.Time

	// workerCompletedMap is a map of Library ids and a boolean to indicate whether the goroutine that was spawned is finished
	workerCompletedMap map[int]bool
}

func (m *Manager) Start(ctx *context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if controller.IsContextFinished(ctx) {
				return
			}
			// Check all Libraries for required scans
			allLibraries := m.ds.Libraries()

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
	discoveredVideos, err := GetVideoFilesFromDir(lib.Folder) // TODO: Abstract this function so it can be mocked out for something else during testing
	if err != nil {
		m.logger.Error(err.Error())
		return
	}

	for _, videoFilepath := range discoveredVideos {
		// Check path against Library path masks
		maskedOut := false
		for _, v := range lib.PathMasks {
			if strings.Contains(videoFilepath, v) {
				m.logger.Debug("%v skipped because of a mask (%v)", videoFilepath, v)
				maskedOut = true
				break
			}
		}
		if maskedOut {
			continue
		}

		if m.ds.IsPathDispatched(videoFilepath) || lib.Queue.InQueuePath(controller.Job{Path: videoFilepath}) {
			continue
		}

		// Read file metadata from a MetadataReader
		fMetadata := m.metadataReader.Read(videoFilepath)

		// Run a CommandDecider against the metadata to determine what FFMpeg command to run
		runCmd, commandSlice := m.commandDecider.Decide(fMetadata)
		if !runCmd {
			m.logger.Debug("Skipping %v because CommandDecider returned a do not run status bool and the following command slice: %v", videoFilepath, commandSlice)
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

		m.ds.SaveLibrary(lib)
	}
}

func (m *Manager) ImportCompletedJobs([]controller.Job) {
	m.logger.Critical("Not implemented")
	// TODO: Implement
}

func (m *Manager) LibrarySettings() (ls []controller.Library) {
	m.logger.Critical("Not implemented")
	// TODO: Implement
	return
}

func (m *Manager) LibraryQueues() (lq []controller.LibraryQueue) {
	m.logger.Critical("Not implemented")
	// TODO: Implement
	return
}

func (m *Manager) PopNewJob() (j controller.Job) {
	m.logger.Critical("Not implemented")
	// TODO: Implement
	return
}

func (m *Manager) UpdateLibrarySettings(map[string]controller.Library) {
	m.logger.Critical("Not implemented")
	// TODO: Implement
}
