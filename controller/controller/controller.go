package controller

import (
	"fmt"
	"sync"
	"time"

	"github.com/BrenekH/project-redcedar-controller/config"
)

// Job represents a job in the RedCedar ecosystem.
type Job struct {
	UUID       string
	Path       string
	Parameters JobParameters
	// RawMediaInfo mediainfo.MediaInfo
}

// JobParameters represents the actions that need to be taken against a job.
type JobParameters struct {
	HEVC        bool // true when the file is not HEVC
	Stereo      bool // true when the file is missing a stereo audio track
	Progressive bool // true when the file is interlaced
}

var controllerConfig *config.ControllerConfiguration

var fileSystemLastCheck time.Time
var healthLastCheck time.Time

// JobQueue is the queue of the jobs
var JobQueue Queue = Queue{sync.Mutex{}, make([]interface{}, 0)}

// RunController is a goroutine compliant way to run the controller.
func RunController(inConfig *config.ControllerConfiguration, stopChan *chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		fmt.Println("stopped")
	}()

	controllerConfig = inConfig
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
		for _, vid := range discoveredVideos {
			// TODO: Run MediaInfo(+ other) checks
			job := Job{UUID: "", Path: vid, Parameters: JobParameters{HEVC: false, Stereo: false, Progressive: false}}
			if !JobQueue.InQueue(job) {
				JobQueue.Push(job)
			}
		}
	}

	if time.Since(healthLastCheck) > time.Duration((*controllerConfig).HealthCheckInterval) {
		healthLastCheck = time.Now()
		// fmt.Println("Doing healthCheck")
		// TODO: Health check
	}
}
