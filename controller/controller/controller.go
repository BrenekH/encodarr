package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/BrenekH/logange"
	"github.com/BrenekH/project-redcedar-controller/config"
	"github.com/BrenekH/project-redcedar-controller/mediainfo"
	"github.com/google/uuid"
)

var logger logange.Logger

func init() {
	logger = logange.NewLogger("controller")
}

// Job represents a job in the RedCedar ecosystem.
type Job struct {
	UUID         string              `json:"uuid"`
	Path         string              `json:"path"`
	Parameters   JobParameters       `json:"parameters"`
	RawMediaInfo mediainfo.MediaInfo `json:"media_info"`
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

// EqualUUID is a custom equality check for the Job type that only checks the UUID parameter
func (j Job) EqualUUID(check Job) bool {
	return j.UUID == check.UUID
}

// DispatchedJob represents a dispatched job in the RedCedar ecosystem.
type DispatchedJob struct {
	Job         Job       `json:"job"`
	RunnerName  string    `json:"runner_name"`
	LastUpdated time.Time `json:"last_updated"`
	Status      JobStatus `json:"status"`
}

// JobStatus represents the status of a dispatched job
type JobStatus struct {
	Stage                       string `json:"stage"`
	Percentage                  string `json:"percentage"`
	JobElapsedTime              string `json:"job_elapsed_time"`
	FPS                         string `json:"fps"`
	StageElapsedTime            string `json:"stage_elapsed_time"`
	StageEstimatedTimeRemaining string `json:"stage_estimated_time_remaining"`
}

// JobRequest represents a request for a job
type JobRequest struct {
	RunnerName    string
	ReturnChannel *chan Job
}

var controllerConfig *config.ControllerConfiguration

var fileSystemLastCheck time.Time
var healthLastCheck time.Time

// JobQueue is the queue of the jobs
var JobQueue Queue = Queue{sync.Mutex{}, make([]Job, 0)}

// DispatchedJobs is a collection for all dispatched jobs
var DispatchedJobs DispatchedContainer = DispatchedContainer{sync.Mutex{}, make([]DispatchedJob, 0)}

// JobRequestChannel is a channel used to send new job requests to the Controller
var JobRequestChannel chan JobRequest = make(chan JobRequest)

// CompletedRequestChannel is a channel used to send job completed requests to the Controller
var CompletedRequestChannel chan JobCompleteRequest = make(chan JobCompleteRequest)

// RunController is a goroutine compliant way to run the controller.
func RunController(inConfig *config.ControllerConfiguration, stopChan *chan interface{}, wg *sync.WaitGroup) {
	wg.Add(1) // This is done in the function rather than outside so that we can easily comment out this function in main.go
	defer wg.Done()
	defer logger.Info("Controller successfully stopped")

	controllerConfig = inConfig

	// Read JSON(Dispatched & History) and apply to containers
	DispatchedJobs = readDispatchedFile()
	HistoryEntries = readHistoryFile()

	// Save if they didn't exist before
	DispatchedJobs.Save()
	HistoryEntries.Save()

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
						var j Job
						for {
							// Pop a job off the Queue
							var err error // err must be defined using var instead of := because j won't be set properly otherwise
							j, err = JobQueue.Pop()
							if err != nil {
								if err == ErrEmptyQueue {
									time.Sleep(time.Duration(int64(0.1 * float64(time.Second)))) // Sleep for 0.1 seconds
									continue
								} else {
									logger.Critical(fmt.Sprintf("Got error while popping from queue: %v", err))
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
								logger.Error(fmt.Sprintf("Unexpected error while stating for file: %v", err))
							}
							time.Sleep(time.Duration(int64(0.1 * float64(time.Second)))) // Sleep for 0.1 seconds
						}

						// Add to dispatched jobs
						DispatchedJobs.Add(DispatchedJob{
							Job:         j,
							RunnerName:  val.RunnerName,
							LastUpdated: time.Now(),
							Status: JobStatus{
								Stage:                       "Waiting to start",
								Percentage:                  "N/A",
								JobElapsedTime:              "N/A",
								FPS:                         "N/A",
								StageElapsedTime:            "N/A",
								StageEstimatedTimeRemaining: "N/A",
							},
						})
						DispatchedJobs.Save()

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
	if time.Since(fileSystemLastCheck) > time.Duration((*controllerConfig).FileSystemCheckInterval) {
		fileSystemLastCheck = time.Now()
		discoveredVideos := GetVideoFilesFromDir((*controllerConfig).SearchDir)
		for _, videoFilepath := range discoveredVideos {
			pathJob := Job{UUID: "", Path: videoFilepath, Parameters: JobParameters{}}

			// Is the file already in the queue or dispatched?
			if JobQueue.InQueuePath(pathJob) || DispatchedJobs.InContainerPath(pathJob) {
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
			logger.Info(fmt.Sprintf("Controller: Added %v to the queue\n", job.Path))
		}
	}
}

func healthCheck() {
	if time.Since(healthLastCheck) > time.Duration((*controllerConfig).HealthCheckInterval) {
		healthLastCheck = time.Now()
		for _, v := range DispatchedJobs.Decontain() {
			if time.Since(v.LastUpdated) > time.Duration((*controllerConfig).HealthCheckTimeout) {
				d, _ := DispatchedJobs.PopByUUID(v.Job.UUID)
				logger.Warn(fmt.Sprintf("Depositing %v back into Job queue because of unresponsive Runner\n", d.Job.Path))
				d.Job.UUID = uuid.New().String()
				JobQueue.Push(d.Job)
				//? Do we follow the python controller and add another "thread-safe" container for timedout jobs or do we return 409 for all requests where the uuid can't be found?
			}
		}
	}
}

func readDispatchedFile() DispatchedContainer {
	// Read/unmarshal json from JSONDir/dispatched_jobs.json
	f, err := os.Open(fmt.Sprintf("%v/dispatched_jobs.json", controllerConfig.JSONDir))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to open dispatched_jobs.json because of error: %v\n", err))
		return DispatchedContainer{sync.Mutex{}, make([]DispatchedJob, 0)}
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to read dispatched_jobs.json because of error: %v\n", err))
		return DispatchedContainer{sync.Mutex{}, make([]DispatchedJob, 0)}
	}

	var readJSON []DispatchedJob
	err = json.Unmarshal(b, &readJSON)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to unmarshal dispatched_jobs.json because of error: %v\n", err))
		return DispatchedContainer{sync.Mutex{}, make([]DispatchedJob, 0)}
	}

	// Add into DispatchedContainer and return
	return DispatchedContainer{sync.Mutex{}, readJSON}
}
