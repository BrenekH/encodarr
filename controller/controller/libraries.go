package controller

import (
	"database/sql"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/BrenekH/encodarr/controller/db/dispatched"
	"github.com/BrenekH/encodarr/controller/db/files"
	"github.com/BrenekH/encodarr/controller/db/libraries"
	"github.com/BrenekH/encodarr/controller/mediainfo"
	"github.com/google/uuid"
)

// The purpose of this file is to hold all code relating to the "bussiness code" of libraries.
// It is not meant to hold any data storage logic, that should all be located in the db/libraries package.

func updateLibraryQueue(l libraries.Library, wg *sync.WaitGroup, completeMap *map[int]bool) {
	wg.Add(1)
	defer wg.Done()
	defer func() { (*completeMap)[l.ID] = true }()

	discoveredVideos := GetVideoFilesFromDir(l.Folder)
	for _, videoFilepath := range discoveredVideos {
		filesEntry := files.File{Path: videoFilepath}
		if err := filesEntry.Get(); err != nil {
			if err != sql.ErrNoRows {
				logger.Warn(err.Error())
				continue
			}
			// File is not in the database yet
			if err = filesEntry.Insert(); err != nil {
				logger.Warn(err.Error())
				continue
			}
		}

		pathJob := dispatched.Job{UUID: "", Path: videoFilepath, Parameters: dispatched.JobParameters{}}

		alreadyDispatched, err := dispatched.PathInDB(pathJob.Path)
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		// Has the file already been dispatched or queued?
		if alreadyDispatched || filesEntry.Queued || l.Queue.InQueuePath(pathJob) {
			logger.Debug(fmt.Sprintf("Skipping %v because it was detected as dispatched or queued", videoFilepath))
			continue
		}

		maskedOut := false
		for _, v := range l.PathMasks {
			if strings.Contains(videoFilepath, v) {
				logger.Debug(fmt.Sprintf("%v skipped because of a mask (%v)", videoFilepath, v))
				maskedOut = true
				break
			}
		}
		if maskedOut {
			continue
		}

		fInfo, err := os.Stat(videoFilepath)
		if err != nil {
			logger.Warn(err.Error())
			continue
		}

		runMediaInfo := true
		// We have to set the mod times to UTC because the db returns a different time zone format than os.Stat()
		if fInfo.ModTime().UTC() == filesEntry.ModTime.UTC() {
			logger.Debug(fmt.Sprintf("Skipping mediainfo on %v because the modtime is the same as the cached version", videoFilepath))
			runMediaInfo = false
		} else {
			logger.Debug(fmt.Sprintf("Adding %v to files table", videoFilepath))
			filesEntry.ModTime = fInfo.ModTime()
			if err = filesEntry.Update(); err != nil {
				logger.Warn(err.Error())
			}
		}

		var mediaInfo mediainfo.MediaInfo
		if runMediaInfo {
			mediaInfo, err := mediainfo.GetMediaInfo(videoFilepath)
			if err != nil {
				logger.Error(fmt.Sprintf("Error getting mediainfo for %v: %v", videoFilepath, err))
				continue
			}
			logger.Trace(fmt.Sprintf("Mediainfo object for %v: %v", videoFilepath, mediaInfo))

			// Save MediaInfo to the database
			filesEntry.MediaInfo = mediaInfo
			if err = filesEntry.Update(); err != nil {
				logger.Warn(err.Error())
			}
		} else {
			mediaInfo = filesEntry.MediaInfo
		}

		// Skips the file if it is not an actual media file
		if !mediaInfo.IsMedia() {
			continue
		}

		// TODO: Change to plugin behavior
		// Is the file HDR?
		if mediaInfo.Video.ColorPrimaries == "BT.2020" {
			continue
		}

		stereoAudioTrackExists := false
		for _, v := range mediaInfo.Audio {
			if v.Channels == "2" {
				stereoAudioTrackExists = true
			}
		}

		isHEVC := mediaInfo.Video.Format == "HEVC"

		if isHEVC && stereoAudioTrackExists {
			continue
		}

		u := uuid.New()
		job := dispatched.Job{
			UUID: u.String(),
			Path: videoFilepath,
			Parameters: dispatched.JobParameters{
				HEVC:   !isHEVC,
				Stereo: !stereoAudioTrackExists,
			},
			RawMediaInfo: mediaInfo,
		}

		logger.Trace(fmt.Sprintf("%v isHEVC=%v stereoAudioTrackExists=%v", videoFilepath, isHEVC, stereoAudioTrackExists))

		l.Queue.Push(job)
		filesEntry.Queued = true
		logger.Info(fmt.Sprintf("Added %v to the queue", job.Path))

		if err = l.Update(); err != nil {
			logger.Error(err.Error())
		}

		if err = filesEntry.Update(); err != nil {
			logger.Error(err.Error())
		}
	}
	logger.Debug(fmt.Sprintf("Finished updating Library %v", l.ID))
}

// isJobAvailable loops through the libraries to identify if any have queued jobs ready
// to be dispatched to Runners.
func isJobAvailable() bool {
	allLibraries, err := libraries.All()
	if err != nil {
		logger.Error(err.Error())
		return false
	}

	for _, v := range allLibraries {
		if len(v.Queue.Items) > 0 {
			return true
		}
	}

	return false
}

// popQueuedJob returns a queued job while also removing it from the queue it was pulled from
func popQueuedJob() (dispatched.Job, error) {
	allLibraries, err := libraries.All()
	if err != nil {
		logger.Error(err.Error())
		return dispatched.Job{}, err
	}

	// Sort libraries by decreasing order so that the libraries with the higher priority number dispatch jobs first.
	sort.Slice(allLibraries, func(i, j int) bool {
		return allLibraries[i].Priority > allLibraries[j].Priority
	})

	for _, v := range allLibraries {
		if len(v.Queue.Items) > 0 {
			item, err := v.Queue.Pop()
			if err != nil {
				logger.Error(err.Error())
				return item, err
			}

			if err = v.Update(); err != nil {
				logger.Error(err.Error())
			}
			// TODO: Update files table to indicate that the job is no longer queued
			return item, nil
		}
	}

	return dispatched.Job{}, fmt.Errorf("no queued jobs were found")
}
