package library

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/BrenekH/encodarr/controller"
)

func NewManager(logger controller.Logger, ds controller.LibraryManagerDataStorer) Manager {
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

	// lastCheckedTimes is a map of Library ids and the last time that they were checked.
	lastCheckedTimes map[int]time.Time

	// workerCompletedMap is a map of Library ids and a boolean to indicate whether the goroutine that was spawned is finished
	workerCompletedMap map[int]bool
}

func (l *Manager) Start(ctx *context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if controller.IsContextFinished(ctx) {
				return
			}
			// Check all Libraries for required scans
			allLibraries := l.ds.Libraries()

			for _, lib := range allLibraries {
				t, ok := l.lastCheckedTimes[lib.ID]
				if !ok {
					l.lastCheckedTimes[lib.ID] = time.Unix(0, 0)
					t = l.lastCheckedTimes[lib.ID]
				}

				previousWorkerFinished, ok := l.workerCompletedMap[lib.ID]
				if !ok {
					l.workerCompletedMap[lib.ID] = true
					previousWorkerFinished = l.workerCompletedMap[lib.ID]
				}

				if time.Since(t) > lib.FsCheckInterval && previousWorkerFinished {
					l.logger.Debug("Initiating library (ID: %v) update", lib.ID)
					l.lastCheckedTimes[lib.ID] = time.Now()
					l.workerCompletedMap[lib.ID] = false

					wg.Add(1)
					go l.updateLibraryQueue(ctx, wg, lib)
				}
			}
			time.Sleep(time.Second)
		}
	}()
}

func (l *Manager) updateLibraryQueue(ctx *context.Context, wg *sync.WaitGroup, lib controller.Library) {
	defer wg.Done()
	defer func() { l.workerCompletedMap[lib.ID] = true }()

	// Locate video files
	discoveredVideos, err := GetVideoFilesFromDir(lib.Folder) // TODO: Abstract this function so it can be mocked out for something else during testing
	if err != nil {
		l.logger.Error(err.Error())
		return
	}

	for _, videoFilepath := range discoveredVideos {
		// Apply Library path masks
		maskedOut := false
		for _, v := range lib.PathMasks {
			if strings.Contains(videoFilepath, v) {
				l.logger.Debug("%v skipped because of a mask (%v)", videoFilepath, v)
				maskedOut = true
				break
			}
		}
		if maskedOut {
			continue
		}

		// Read file metadata from a MetadataReader
		// Run a CommandDecider against the metadata to determine what FFMpeg command to run
		// Save to Library queue
	}
}

func (l *Manager) ImportCompletedJobs([]controller.Job) {
	l.logger.Critical("Not implemented")
	// TODO: Implement
}

func (l *Manager) LibrarySettings() (ls []controller.Library) {
	l.logger.Critical("Not implemented")
	// TODO: Implement
	return
}

func (l *Manager) LibraryQueues() (lq []controller.LibraryQueue) {
	l.logger.Critical("Not implemented")
	// TODO: Implement
	return
}

func (l *Manager) PopNewJob() (j controller.Job) {
	l.logger.Critical("Not implemented")
	// TODO: Implement
	return
}

func (l *Manager) UpdateLibrarySettings(map[string]controller.Library) {
	l.logger.Critical("Not implemented")
	// TODO: Implement
}
