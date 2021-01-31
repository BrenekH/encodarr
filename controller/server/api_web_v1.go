package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/BrenekH/project-redcedar-controller/controller"
)

type queueJSONResponse struct {
	JobQueue []controller.Job `json:"queue"`
}

type runningJSONResponse struct {
	DispatchedJobs []filteredDispatchedJob `json:"jobs"`
}

type filteredDispatchedJob struct {
	Job        filteredJob          `json:"job"`
	RunnerName string               `json:"runner_name"`
	Status     controller.JobStatus `json:"status"`
}

type filteredJob struct {
	UUID       string                   `json:"uuid"`
	Path       string                   `json:"path"`
	Parameters controller.JobParameters `json:"parameters"`
}

func makeFilteredDispatchedJobs() runningJSONResponse {
	dispatchedJobsSlice := controller.DispatchedJobs.Decontain()
	runningJSONResponseStruct := runningJSONResponse{DispatchedJobs: make([]filteredDispatchedJob, len(dispatchedJobsSlice))}

	for i, dJob := range dispatchedJobsSlice {
		runningJSONResponseStruct.DispatchedJobs[i] = filteredDispatchedJob{
			Job: filteredJob{
				UUID:       dJob.Job.UUID,
				Path:       dJob.Job.Path,
				Parameters: dJob.Job.Parameters,
			},
			RunnerName: dJob.RunnerName,
			Status:     dJob.Status,
		}
	}

	return runningJSONResponseStruct
}

// getRunning is a HTTP handler that returns the current running jobs in a JSON response.
func getRunning(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		runningJSONBytes, err := json.Marshal(makeFilteredDispatchedJobs())
		if err != nil {
			serverError(w, r, fmt.Sprintf("Error marshaling Job queue to json: %v", err))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(runningJSONBytes)
	default:
		methodForbidden(w, r)
	}
}

// getQueue is a HTTP handler that returns the current queue in a JSON response.
func getQueue(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		jsonResponseStruct := queueJSONResponse{controller.JobQueue.Dequeue()}
		queueJSONBytes, err := json.Marshal(jsonResponseStruct)
		if err != nil {
			serverError(w, r, fmt.Sprintf("Error marshaling Job queue to json: %v", err))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(queueJSONBytes)
	default:
		methodForbidden(w, r)
	}
}

// TODO: Complete GET history
// getHistory is a HTTP handler that returns the current history in a JSON response.
func getHistory(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"test": true}`))
	default:
		methodForbidden(w, r)
	}
}

func registerWebAPIv1Handlers() {
	r := newSubRouter("/api/web/v1")

	r.HandleFunc("/running", getRunning)
	r.HandleFunc("/queue", getQueue)
	r.HandleFunc("/history", getHistory)
}
