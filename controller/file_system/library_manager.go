package file_system

import (
	"context"
	"sync"
	"time"

	"github.com/BrenekH/encodarr/controller"
)

func NewLibraryManager(logger controller.Logger) LibraryManager {
	return LibraryManager{
		logger: logger,

		lastCheckedTimes:   make(map[int]time.Time),
		workerCompletedMap: make(map[int]bool),
	}
}

type LibraryManager struct {
	logger controller.Logger
	ds     controller.LibraryManagerDataStorer

	// lastCheckedTimes is a map of Library ids and the last time that they were checked.
	lastCheckedTimes map[int]time.Time

	// workerCompletedMap is a map of Library ids and a boolean to indicate whether the goroutine that was spawned is finished
	workerCompletedMap map[int]bool
}

func (l *LibraryManager) Start(ctx *context.Context, wg *sync.WaitGroup) {
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

func (l *LibraryManager) updateLibraryQueue(ctx *context.Context, wg *sync.WaitGroup, lib controller.Library) {
	defer wg.Done()
	l.logger.Critical("Not implemented")
	// Locate media files
	// Read file metadata from a MetadataReader
	// Cache the file metadata using a DataStorer (caching could be integrated into the MetadataReader)
	// Run a CommandDecider against the metadata to determine what FFMpeg command to run
	// Save to Library queue
}

func (l *LibraryManager) ImportCompletedJobs([]controller.Job) {
	l.logger.Critical("Not implemented")
}

func (l *LibraryManager) LibrarySettings() (ls []controller.Library) {
	l.logger.Critical("Not implemented")
	return
}

func (l *LibraryManager) LibraryQueues() (lq []controller.LibraryQueue) {
	l.logger.Critical("Not implemented")
	return
}

func (l *LibraryManager) PopNewJob() (j controller.Job) {
	l.logger.Critical("Not implemented")
	return
}

func (l *LibraryManager) UpdateLibrarySettings(map[string]controller.Library) {
	l.logger.Critical("Not implemented")
}
