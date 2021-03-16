package controller

import (
	"fmt"
	"strings"

	"github.com/BrenekH/project-redcedar-controller/db/dispatched"
	"github.com/BrenekH/project-redcedar-controller/db/libraries"
	"github.com/BrenekH/project-redcedar-controller/mediainfo"
	"github.com/google/uuid"
)

// The purpose of this file is to hold all code relating to the "bussiness code" of libraries.
// It is not meant to hold any data storage logic, that should all be located in the db/libraries package.

func updateLibraryQueue(l libraries.Library) {
	discoveredVideos := GetVideoFilesFromDir(l.Folder)
	for _, videoFilepath := range discoveredVideos {
		// TODO: Check modtime against file_cache table

		pathJob := dispatched.Job{UUID: "", Path: videoFilepath, Parameters: dispatched.JobParameters{}}

		// TODO: Change to checking the file_cache table and dispatched jobs
		// Is the file already in the queue or dispatched?
		alreadyInDB, err := dispatched.PathInDB(pathJob.Path)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		if l.Queue.InQueuePath(pathJob) || alreadyInDB {
			continue
		}

		// TODO: Change to path mask
		// Is the file 'optimized' by Plex?
		if strings.Contains(videoFilepath, "Plex Versions") {
			continue
		}

		mediainfo, err := mediainfo.GetMediaInfo(videoFilepath)
		if err != nil {
			logger.Error(fmt.Sprintf("Error getting mediainfo for %v: %v", videoFilepath, err))
			continue
		}
		logger.Trace(fmt.Sprintf("Mediainfo object for %v: %v", videoFilepath, mediainfo))

		// Skips the file if it is not an actual media file
		if !mediainfo.IsMedia() {
			continue
		}

		// TODO: Change to plugin behavior
		// Is the file HDR?
		if mediainfo.Video.ColorPrimaries == "BT.2020" {
			continue
		}

		stereoAudioTrackExists := false
		for _, v := range mediainfo.Audio {
			if v.Channels == "2" {
				stereoAudioTrackExists = true
			}
		}

		isHEVC := mediainfo.Video.Format == "HEVC"

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
			RawMediaInfo: mediainfo,
		}

		logger.Trace(fmt.Sprintf("%v isHEVC=%v stereoAudioTrackExists=%v", videoFilepath, isHEVC, stereoAudioTrackExists))

		l.Queue.Push(job)
		logger.Info(fmt.Sprintf("Added %v to the queue", job.Path))
	}
}
