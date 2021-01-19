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
	HEVC        bool `json:"hevc"`        // true when the file is not HEVC
	Stereo      bool `json:"stereo"`      // true when the file is missing a stereo audio track
	Progressive bool `json:"progressive"` // true when the file is interlaced
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

// EqualPath is a custom equality check for the Job type that ignores the UUID but checks everything else
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
			controllerLoop()
		case <-*stopChan:
			return
		}
	}
}

func controllerLoop() {
	if time.Since(fileSystemLastCheck) > time.Duration((*controllerConfig).FileSystemCheckInterval) {
		fileSystemLastCheck = time.Now()
		// fmt.Println("Doing fileSystemCheck")
		// TODO: File system check
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

			// TODO: Set platform binary in mediainfo.go, not here
			windowsMediaInfo := "MediaInfo.exe"
			err := mediainfo.SetMediaInfoBinary(windowsMediaInfo)
			if err != nil {
				log.Fatal(err)
			}

			mediainfo, err := mediainfo.GetMediaInfo(videoFilepath)
			if err != nil {
				log.Fatal(err)
			}

			// Is the file HDR?
			if mediainfo.Video.ColorPrimaries == "BT.2020" {
				continue
			}

			u := uuid.New()
			job := Job{
				UUID: u.String(),
				Path: videoFilepath,
				Parameters: JobParameters{
					HEVC:        false,
					Stereo:      false,
					Progressive: false,
				},
				RawMediaInfo: mediainfo,
			}

			JobQueue.Push(job)
			log.Printf("Controller: Added %v to the queue\n", job.Path)
		}
	}

	if time.Since(healthLastCheck) > time.Duration((*controllerConfig).HealthCheckInterval) {
		healthLastCheck = time.Now()
		// fmt.Println("Doing healthCheck")
		// TODO: Health check
	}
}
