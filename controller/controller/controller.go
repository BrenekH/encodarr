package controller

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/BrenekH/logange"
	"github.com/BrenekH/project-redcedar-controller/config"
	"github.com/BrenekH/project-redcedar-controller/db/dispatched"
	"github.com/BrenekH/project-redcedar-controller/mediainfo"
	"github.com/BrenekH/project-redcedar-controller/options"
	"github.com/google/uuid"
)

var logger logange.Logger

func init() {
	logger = logange.NewLogger("controller")
}

// DispatchedJob represents a dispatched job in the RedCedar ecosystem.
type DispatchedJob struct {
	Job         dispatched.Job       `json:"job"`
	RunnerName  string               `json:"runner_name"`
	LastUpdated time.Time            `json:"last_updated"`
	Status      dispatched.JobStatus `json:"status"`
}

// JobRequest represents a request for a job
type JobRequest struct {
	RunnerName    string
	ReturnChannel *chan dispatched.Job
}

var fileSystemLastCheck time.Time
var healthLastCheck time.Time

// JobQueue is the queue of the jobs
var JobQueue Queue = Queue{sync.Mutex{}, make([]dispatched.Job, 0)}

// JobRequestChannel is a channel used to send new job requests to the Controller
var JobRequestChannel chan JobRequest = make(chan JobRequest)

// CompletedRequestChannel is a channel used to send job completed requests to the Controller
var CompletedRequestChannel chan JobCompleteRequest = make(chan JobCompleteRequest)

// RunController is a goroutine compliant way to run the controller.
func RunController(stopChan *chan interface{}, wg *sync.WaitGroup) {
	wg.Add(1) // This is done in the function rather than outside so that we can easily comment out this function in main.go
	defer wg.Done()
	defer logger.Info("Controller successfully stopped")

	// Start the job request handler
	go jobRequestHandler(&JobRequestChannel, stopChan, wg)

	// Start the completed request handler
	go completedLooper(&CompletedRequestChannel, stopChan, wg)

	// This loop is in charge of running the controller logic until the stop signal channel "stopChan" has a value pushed to it
	for {
		select {
		default:
			fileSystemCheck()
			healthCheck()
		case <-*stopChan:
			return
		}
		time.Sleep(time.Duration(int64(0.1 * float64(time.Second)))) // Sleep for 0.1 seconds
	}
}

// jobRequestHandler continuously checks the requestChan interface and responds with a job
func jobRequestHandler(requestChan *chan JobRequest, stopChan *chan interface{}, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	defer logger.Info("jobRequestHandler stopped")

	for {
		select {
		default:
			if !JobQueue.Empty() {
				select {
				case val, ok := <-*requestChan:
					if ok {
						var j dispatched.Job
						for {
							// Pop a job off the Queue
							var err error // err must be defined using var instead of := because j won't be set properly otherwise
							j, err = JobQueue.Pop()
							if err != nil {
								if err == ErrEmptyQueue {
									time.Sleep(time.Duration(int64(0.1 * float64(time.Second)))) // Sleep for 0.1 seconds
									continue
								} else {
									logger.Critical(fmt.Sprintf("Got unexpected error while popping from queue: %v", err))
								}
							}

							// Check if the job is still valid
							if _, err := os.Stat(j.Path); err == nil {
								// TODO: Do more than just check if it exists (verify hevc and stereo attributes)
								break
							} else if os.IsNotExist(err) {
								// File does not exist. Do not add back into queue
								continue
							} else {
								// File may or may not exist. Error has more details.
								logger.Error(fmt.Sprintf("Unexpected error while stating file '%v': %v", j.Path, err))
							}
							time.Sleep(time.Duration(int64(0.1 * float64(time.Second)))) // Sleep for 0.1 seconds
						}

						// Add to dispatched jobs
						dJob := dispatched.DJob{
							UUID:        j.UUID,
							Job:         j,
							Runner:      val.RunnerName,
							LastUpdated: time.Now(),
							Status: dispatched.JobStatus{
								Stage:                       "Copying to Runner",
								Percentage:                  "0",
								JobElapsedTime:              "N/A",
								FPS:                         "N/A",
								StageElapsedTime:            "N/A",
								StageEstimatedTimeRemaining: "N/A",
							},
						}
						err := dJob.Insert()
						if err != nil {
							logger.Error(fmt.Sprintf("Error saving dispatched jobs: %v", err.Error()))
						}

						// Return Job struct in return channel
						*val.ReturnChannel <- j
					} else {
						// Channel closed. Stop handler.
						return
					}
				case <-*stopChan:
					return
				default:
				}
			}
		case <-*stopChan:
			return
		}
		time.Sleep(time.Duration(int64(0.1 * float64(time.Second)))) // Sleep for 0.1 seconds
	}
}

func fileSystemCheck() {
	if time.Since(fileSystemLastCheck) > time.Duration(config.Global.FileSystemCheckInterval) {
		logger.Debug("Starting file system check")
		fileSystemLastCheck = time.Now()
		discoveredVideos := GetVideoFilesFromDir(options.SearchDir())
		for _, videoFilepath := range discoveredVideos {
			pathJob := dispatched.Job{UUID: "", Path: videoFilepath, Parameters: dispatched.JobParameters{}}

			// Is the file already in the queue or dispatched?
			alreadyInDB, err := dispatched.PathInDB(pathJob.Path)
			if err != nil {
				logger.Error(err.Error())
				continue
			}

			if JobQueue.InQueuePath(pathJob) || alreadyInDB {
				continue
			}

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

			JobQueue.Push(job)
			logger.Info(fmt.Sprintf("Added %v to the queue", job.Path))
		}
		logger.Debug("File system check complete")
	}
}

func healthCheck() {
	if time.Since(healthLastCheck) > time.Duration(config.Global.HealthCheckInterval) {
		healthLastCheck = time.Now()
		logger.Debug("Starting health check")
		dJobs, err := dispatched.All()
		if err != nil {
			logger.Error(err.Error())
			return
		}

		for _, v := range dJobs {
			if time.Since(v.LastUpdated) > time.Duration(config.Global.HealthCheckTimeout) {
				d := dispatched.DJob{UUID: v.Job.UUID}
				if err = d.Get(); err != nil {
					logger.Error(err.Error())
					continue
				}

				logger.Warn(fmt.Sprintf("Removing %v from dispatched jobs because of unresponsive Runner", d.Job.Path))

				if err = d.Delete(); err != nil {
					logger.Error(err.Error())
					continue
				}
				// TODO: Add back into library queue
			}
		}
		logger.Debug("Health check complete")
	}
}
