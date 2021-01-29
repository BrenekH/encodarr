package controller

import (
	"log"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/BrenekH/project-redcedar-controller/config"
	"github.com/BrenekH/project-redcedar-controller/mediainfo"
	"github.com/google/uuid"
)

// Job represents a job in the RedCedar ecosystem.
type Job struct {
	UUID         string              `json:"uuid"`
	Path         string              `json:"path"`
	Parameters   JobParameters       `json:"parameters"`
	RawMediaInfo mediainfo.MediaInfo `json:"raw_media_info"`
}

// JobParameters represents the actions that need to be taken against a job.
type JobParameters struct {
	HEVC   bool `json:"hevc"`   // true when the file is not HEVC
	Stereo bool `json:"stereo"` // true when the file is missing a stereo audio track
}

// Equal is a custom equality check for the Job type
func (j Job) Equal(check Job) bool {
	if j.UUID != check.UUID {
		return false
	}
	if j.Path != check.Path {
		return false
	}
	if !reflect.DeepEqual(j.Parameters, check.Parameters) {
		return false
	}
	return true
}

// EqualPath is a custom equality check for the Job type that only checks the Path parameter
func (j Job) EqualPath(check Job) bool {
	return j.Path == check.Path
}

var controllerConfig *config.ControllerConfiguration

var fileSystemLastCheck time.Time
var healthLastCheck time.Time

// JobQueue is the queue of the jobs
var JobQueue Queue = Queue{sync.Mutex{}, make([]Job, 0)}

// RunController is a goroutine compliant way to run the controller.
func RunController(inConfig *config.ControllerConfiguration, stopChan *chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer log.Println("Controller: Successfully stopped")

	controllerConfig = inConfig

	// This loop is in charge of running the controller logic until the stop signal channel "stopChan" has a value pushed to it
	for {
		select {
		default:
			fileSystemCheck()
			healthCheck()
		case <-*stopChan:
			return
		}
	}
}

func fileSystemCheck() {
	if time.Since(fileSystemLastCheck) > time.Duration((*controllerConfig).FileSystemCheckInterval) {
		fileSystemLastCheck = time.Now()
		discoveredVideos := GetVideoFilesFromDir((*controllerConfig).SearchDir)
		for _, videoFilepath := range discoveredVideos {
			pathJob := Job{UUID: "", Path: videoFilepath, Parameters: JobParameters{}}

			// Is the file already in the queue?
			if JobQueue.InQueuePath(pathJob) {
				continue
			}

			// Is the file 'optimized' by Plex?
			if strings.Contains(videoFilepath, "Plex Versions") {
				continue
			}

			mediainfo, err := mediainfo.GetMediaInfo(videoFilepath)
			if err != nil {
				log.Fatal(err)
			}

			// Skips the file if it is not an actual media file
			if !mediainfo.IsMedia() {
				continue
			}

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
			job := Job{
				UUID: u.String(),
				Path: videoFilepath,
				Parameters: JobParameters{
					HEVC:   !isHEVC,
					Stereo: !stereoAudioTrackExists,
				},
				RawMediaInfo: mediainfo,
			}

			JobQueue.Push(job)
			log.Printf("Controller: Added %v to the queue\n", job.Path)
		}
	}
}

func healthCheck() {
	if time.Since(healthLastCheck) > time.Duration((*controllerConfig).HealthCheckInterval) {
		healthLastCheck = time.Now()
		// TODO: Health check
	}
}
