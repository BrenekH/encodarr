package controller

import (
	"fmt"
	"sync"
	"time"

	"github.com/BrenekH/project-redcedar-controller/config"
	"github.com/BrenekH/project-redcedar-controller/mediainfo"
)

// Job represents a job in the RedCedar ecosystem.
type Job struct {
	UUID         string
	Path         string
	Parameters   JobParameters
	RawMediaInfo mediainfo.MediaInfo
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
	// TODO: Check if file system check should be run and complete it if necessary
	if time.Since(fileSystemLastCheck) > time.Duration((*controllerConfig).FileSystemCheckInterval) {
		fileSystemLastCheck = time.Now()
		fmt.Println("Doing fileSystemCheck")
	}
	// TODO: Check if health check should be run and complete it if necessary
	if time.Since(healthLastCheck) > time.Duration((*controllerConfig).HealthCheckInterval) {
		healthLastCheck = time.Now()
		fmt.Println("Doing healthCheck")
	}
}
