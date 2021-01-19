package controller

import (
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/BrenekH/project-redcedar-controller/config"
	"github.com/google/uuid"
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

// Equal is a custom equality check for the Job type that ignores the UUID but checks everything else
func (j Job) Equal(check Job) bool {
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

	// This loop is in charge of running the controller logic until the stop signal channel stopChan has a value pushed to it
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
			// TODO: Run MediaInfo(+ other) checks
			u := uuid.New()
			job := Job{UUID: u.String(), Path: videoFilepath, Parameters: JobParameters{HEVC: false, Stereo: false, Progressive: false}}
			if !JobQueue.InQueuePath(job) {
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
