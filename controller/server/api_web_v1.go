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

type transformedHistoryEntry struct {
	File              string   `json:"file"`
	DateTimeCompleted string   `json:"datetime_completed"`
	Warnings          []string `json:"warnings"`
	Errors            []string `json:"errors"`
}

type historyJSON struct {
	History []transformedHistoryEntry `json:"history"`
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

// getHistory is a HTTP handler that returns the current history in a JSON response.
func getHistory(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Get slice of HistoryEntries (Decontain)
		hE := controller.HistoryEntries.Decontain()
		hELen := len(hE)
		h := make([]transformedHistoryEntry, hELen)

		// Change datetime into human-readable format
		for i, v := range hE {
			dt := v.DateTimeCompleted
			h[hELen-i] = transformedHistoryEntry{
				File: v.File,
				DateTimeCompleted: fmt.Sprintf("%02d-%02d-%d %02d:%02d:%02d",
					dt.Month(), dt.Day(), dt.Year(),
					dt.Hour(), dt.Minute(), dt.Second()),
				Warnings: v.Warnings,
				Errors:   v.Errors,
			}
		}

		// Send JSON to client
		historyJSONBytes, err := json.Marshal(historyJSON{h})
		if err != nil {
			serverError(w, r, fmt.Sprintf("Error marshaling Job history to json: %v", err))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(historyJSONBytes)
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
