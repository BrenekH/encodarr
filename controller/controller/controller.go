package controller

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/BrenekH/logange"
	"github.com/BrenekH/project-redcedar-controller/config"
	"github.com/BrenekH/project-redcedar-controller/db/dispatched"
	"github.com/BrenekH/project-redcedar-controller/db/libraries"
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

// fsCheckTimes is a map of Library ids and the last time that they were checked.
var fsCheckTimes map[int]time.Time = make(map[int]time.Time)

// healthLastCheck holds the last time a health check was performed.
var healthLastCheck time.Time

// JobQueue is the queue of the jobs
var JobQueue libraries.Queue = libraries.Queue{Items: make([]dispatched.Job, 0)}

// JobRequestChannel is a channel used to send new job requests to the Controller
var JobRequestChannel chan JobRequest = make(chan JobRequest)

// CompletedRequestChannel is a channel used to send job completed requests to the Controller
var CompletedRequestChannel chan JobCompleteRequest = make(chan JobCompleteRequest)

// jobRequests holds all of the requests until they can be resolved
var jobRequests []JobRequest = make([]JobRequest, 0)

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
			fileSystemCheck(wg)
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
			select {
			case c := <-*requestChan:
				jobRequests = append(jobRequests, c)
			default:
			}

			if len(jobRequests) != 0 {
				if isJobAvailable() {
					jR, err := popJobRequest()
					if err != nil {
						logger.Warn(err.Error())
						continue
					}

					var j dispatched.Job
					doNotUseJob := false
					for {
						// Pop a job off the Queue
						var err error // err must be defined using var instead of := because j won't be set properly otherwise
						j, err = popQueuedJob()
						if err != nil {
							logger.Error(err.Error())
							doNotUseJob = true
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
					if doNotUseJob {
						continue
					}

					// Add to dispatched jobs
					dJob := dispatched.DJob{
						UUID:        j.UUID,
						Job:         j,
						Runner:      jR.RunnerName,
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
					err = dJob.Insert()
					if err != nil {
						logger.Error(fmt.Sprintf("Error saving dispatched jobs: %v", err.Error()))
					}

					// Return Job struct in return channel
					*jR.ReturnChannel <- j
				}
			}
		case <-*stopChan:
			return
		}
		time.Sleep(time.Duration(int64(0.1 * float64(time.Second)))) // Sleep for 0.1 seconds
	}
}

func fileSystemCheck(wg *sync.WaitGroup) {
	allLibraries, err := libraries.All()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	for _, l := range allLibraries {
		t, ok := fsCheckTimes[l.ID]
		if !ok {
			fsCheckTimes[l.ID] = time.Unix(0, 0)
			t = fsCheckTimes[l.ID]
		}

		if time.Since(t) > l.FsCheckInterval {
			logger.Debug(fmt.Sprintf("Initiating library (ID: %v) update", l.ID))
			fsCheckTimes[l.ID] = time.Now()
			go updateLibraryQueue(l, wg)
		}
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

// popJobRequest returns the first element of the jobRequests slice
// and shifts the remaining items up one slot.
func popJobRequest() (JobRequest, error) {
	if len(jobRequests) == 0 {
		return JobRequest{}, fmt.Errorf("jobRequests is empty")
	}
	item := jobRequests[0]
	jobRequests[0] = JobRequest{} // Hopefully this garbage collects properly
	jobRequests = jobRequests[1:]
	return item, nil
}
