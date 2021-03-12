package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/BrenekH/project-redcedar-controller/config"
	"github.com/BrenekH/project-redcedar-controller/controller"
	"github.com/BrenekH/project-redcedar-controller/db/dispatched"
	"github.com/BrenekH/project-redcedar-controller/db/history"
)

type queueJSONResponse struct {
	JobQueue []dispatched.Job `json:"queue"`
}

type runningJSONResponse struct {
	DispatchedJobs []filteredDispatchedJob `json:"jobs"`
}

type filteredDispatchedJob struct {
	Job        filteredJob          `json:"job"`
	RunnerName string               `json:"runner_name"`
	Status     dispatched.JobStatus `json:"status"`
}

type filteredJob struct {
	UUID       string                   `json:"uuid"`
	Path       string                   `json:"path"`
	Parameters dispatched.JobParameters `json:"parameters"`
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

type settingsJSON struct {
	FileSystemCheckInterval string
	HealthCheckInterval     string
	HealthCheckTimeout      string
	LogVerbosity            string
	SmallerFiles            bool
}

func makeFilteredDispatchedJobs() runningJSONResponse {
	dispatchedJobsSlice, _ := dispatched.All()
	runningJSONResponseStruct := runningJSONResponse{DispatchedJobs: make([]filteredDispatchedJob, len(dispatchedJobsSlice))}

	for i, dJob := range dispatchedJobsSlice {
		runningJSONResponseStruct.DispatchedJobs[i] = filteredDispatchedJob{
			Job: filteredJob{
				UUID:       dJob.Job.UUID,
				Path:       dJob.Job.Path,
				Parameters: dJob.Job.Parameters,
			},
			RunnerName: dJob.Runner,
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
		hE, err := history.All()
		if err != nil {
			logger.Error(err.Error())
			serverError(w, r, err.Error())
			return
		}

		h := make([]transformedHistoryEntry, len(hE))

		// Change datetime into human-readable format
		for i, v := range hE {
			dt := v.DateTimeCompleted
			h[i] = transformedHistoryEntry{
				File: v.Filename,
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

func settings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rS := settingsJSON{
			FileSystemCheckInterval: time.Duration(config.Global.FileSystemCheckInterval).String(),
			HealthCheckInterval:     time.Duration(config.Global.HealthCheckInterval).String(),
			HealthCheckTimeout:      time.Duration(config.Global.HealthCheckTimeout).String(),
			LogVerbosity:            config.RootFileHandler.LevelString(),
			SmallerFiles:            config.Global.SmallerFiles,
		}
		b, err := json.Marshal(rS)
		if err != nil {
			serverError(w, r, fmt.Sprintf("Failed to marshal settingsJSON: %v", err))
		}
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	case http.MethodPut:
		b, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to read request body: %v", err))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(""))
			return
		}

		rS := settingsJSON{
			FileSystemCheckInterval: time.Duration(config.Global.FileSystemCheckInterval).String(),
			HealthCheckInterval:     time.Duration(config.Global.HealthCheckInterval).String(),
			HealthCheckTimeout:      time.Duration(config.Global.HealthCheckTimeout).String(),
			LogVerbosity:            config.RootFileHandler.LevelString(),
			SmallerFiles:            config.Global.SmallerFiles,
		}
		err = json.Unmarshal(b, &rS)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to unmarshal settings put request body: %v", err))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Malformed request body"))
			return
		}

		td, err := time.ParseDuration(rS.FileSystemCheckInterval)
		if err == nil {
			config.Global.FileSystemCheckInterval = int(td)
		}

		td, err = time.ParseDuration(rS.HealthCheckInterval)
		if err == nil {
			config.Global.HealthCheckInterval = int(td)
		}

		td, err = time.ParseDuration(rS.HealthCheckTimeout)
		if err == nil {
			config.Global.HealthCheckTimeout = int(td)
		}

		err = config.SetRootFHVerbosity(rS.LogVerbosity)
		if err != nil {
			logger.Warn(err.Error())
		} else {
			config.Global.LogVerbosity = rS.LogVerbosity
		}

		config.Global.SmallerFiles = rS.SmallerFiles

		err = config.SaveGlobal()
		if err != nil {
			logger.Error(err.Error())
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(""))
	default:
		methodForbidden(w, r)
	}
}

func registerWebAPIv1Handlers() {
	r := newSubRouter("/api/web/v1")

	r.HandleFunc("/running", getRunning)
	r.HandleFunc("/queue", getQueue)
	r.HandleFunc("/history", getHistory)
	r.HandleFunc("/settings", settings)
}
