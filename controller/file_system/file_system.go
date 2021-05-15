package file_system

import (
	"context"

	"github.com/BrenekH/encodarr/controller"
)

func NewLibraryManager(logger controller.Logger) LibraryManager {
	return LibraryManager{
		logger: logger,
	}
}

type LibraryManager struct {
	logger controller.Logger
}

func (l *LibraryManager) Start(ctx *context.Context) {
	l.logger.Critical("Not implemented")
	// Check all Libraries for required scans
	// Scan goroutine
	//   - Locate media files
	//   - Read file metadata from a MetadataReader
	//   - Cache the file metadata using a DataStorer (caching could be integrated into the MetadataReader)
	//   - Run a CommandDecider against the metadata to determine what FFMpeg command to run
	//   - Save to Library queue
}

func (l *LibraryManager) ImportCompletedJobs([]controller.Job) {
	l.logger.Critical("Not implemented")
}

func (l *LibraryManager) LibrarySettings() (ls []controller.LibrarySettings) {
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

func (l *LibraryManager) UpdateLibrarySettings(map[string]controller.LibrarySettings) {
	l.logger.Critical("Not implemented")
}
