package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/BrenekH/encodarr/controller/config"
	"github.com/BrenekH/encodarr/controller/controller"
	"github.com/BrenekH/encodarr/controller/db/dispatched"
	"github.com/BrenekH/encodarr/controller/db/history"
	"github.com/BrenekH/encodarr/controller/db/libraries"
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

type waitingRunners struct {
	Runners []string
}

type libraryJSON struct {
	ID              int                      `json:"id"`
	Folder          string                   `json:"folder"`
	Priority        int                      `json:"priority"`
	FsCheckInterval string                   `json:"fs_check_interval"`
	Pipeline        libraries.PluginPipeline `json:"pipeline"`
	Queue           libraries.Queue          `json:"queue"`
	PathMasks       []string                 `json:"path_masks"`
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
		jsonResponseStruct := queueJSONResponse{JobQueue: make([]dispatched.Job, 0)}

		allLibraries, err := libraries.All()
		if err != nil {
			logger.Error(err.Error())
			return
		}

		for _, v := range allLibraries {
			jsonResponseStruct.JobQueue = append(jsonResponseStruct.JobQueue, v.Queue.Items...)
		}

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
			HealthCheckInterval: time.Duration(config.Global.HealthCheckInterval).String(),
			HealthCheckTimeout:  time.Duration(config.Global.HealthCheckTimeout).String(),
			LogVerbosity:        config.RootFileHandler.LevelString(),
			SmallerFiles:        config.Global.SmallerFiles,
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
			HealthCheckInterval: time.Duration(config.Global.HealthCheckInterval).String(),
			HealthCheckTimeout:  time.Duration(config.Global.HealthCheckTimeout).String(),
			LogVerbosity:        config.RootFileHandler.LevelString(),
			SmallerFiles:        config.Global.SmallerFiles,
		}
		err = json.Unmarshal(b, &rS)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to unmarshal settings put request body: %v", err))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Malformed request body"))
			return
		}

		td, err := time.ParseDuration(rS.HealthCheckInterval)
		if err == nil {
			config.Global.HealthCheckInterval = uint64(td)
		}

		td, err = time.ParseDuration(rS.HealthCheckTimeout)
		if err == nil {
			config.Global.HealthCheckTimeout = uint64(td)
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

// getWaitingRunners is a HTTP handler that returns all runners waiting for a job in a JSON response.
func getWaitingRunners(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		runners := make([]string, len(controller.JobRequests))

		for _, v := range controller.JobRequests {
			runners = append(runners, v.RunnerName)
		}

		if len(runners) > 0 {
			runners = runners[1:]
		}
		wR := waitingRunners{Runners: runners}

		b, err := json.Marshal(wR)
		if err != nil {
			logger.Error(err.Error())
			serverError(w, r, err.Error())
			return
		}

		w.Write(b)
	default:
		methodForbidden(w, r)
	}
}

// getAllLibraryIDs is a HTTP handler that returns all of the library's IDs
func getAllLibraryIDs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		allLibs, err := libraries.All()
		if err != nil {
			logger.Error(err.Error())
			serverError(w, r, err.Error())
			return
		}

		ids := make([]int, len(allLibs))
		for k, v := range allLibs {
			ids[k] = v.ID
		}

		b, err := json.Marshal(struct{ IDs []int }{ids})
		if err != nil {
			logger.Error(err.Error())
			serverError(w, r, err.Error())
			return
		}
		w.Write(b)
	default:
		methodForbidden(w, r)
	}
}

// handleLibrary is a HTTP handler than takes care of the management of a Library
func handleLibrary(w http.ResponseWriter, r *http.Request) {
	libraryID := r.URL.Path[len("/api/web/v1/library/"):]

	if libraryID == "new" && r.Method == http.MethodPost {
		readBytes, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error(err.Error())
			serverError(w, r, err.Error())
			return
		}

		interimNewLib := libraryJSON{}
		err = json.Unmarshal(readBytes, &interimNewLib)
		if err != nil {
			logger.Error(err.Error())
			serverError(w, r, err.Error())
			return
		}

		newLib := libraries.Library{
			Folder:    interimNewLib.Folder,
			Priority:  interimNewLib.Priority,
			Pipeline:  interimNewLib.Pipeline,
			PathMasks: interimNewLib.PathMasks,
		}

		td, err := time.ParseDuration(interimNewLib.FsCheckInterval)
		if err == nil {
			newLib.FsCheckInterval = td
		}

		if err = newLib.Create(); err != nil {
			logger.Error(err.Error())
			serverError(w, r, err.Error())
			return
		}

		w.Write([]byte(""))
		return
	}

	// Transform the string libraryID into an int intLibID
	temp, err := strconv.ParseInt(libraryID, 0, 0)
	if err != nil {
		logger.Error(err.Error())
		serverError(w, r, err.Error())
		return
	}
	intLibID := int(temp)

	// Validate libraryID
	lib := libraries.Library{ID: intLibID}
	if err = lib.Get(); err != nil {
		logger.Error(err.Error())
		serverError(w, r, err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		toSend := libraryJSON{lib.ID, lib.Folder, lib.Priority, lib.FsCheckInterval.String(), lib.Pipeline, lib.Queue, lib.PathMasks}
		b, err := json.Marshal(toSend)
		if err != nil {
			logger.Error(err.Error())
			serverError(w, r, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	case http.MethodPut:
		// Technically, there is a security flaw where an attacker can set the id in their request
		// to a different library and overwrite a different library, but it's not like this API is locked down at all
		// so does it really matter?
		readBytes, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error(err.Error())
			serverError(w, r, err.Error())
			return
		}

		uLib := libraryJSON{}
		err = json.Unmarshal(readBytes, &uLib)
		if err != nil {
			logger.Error(err.Error())
			serverError(w, r, err.Error())
			return
		}

		lib.Folder = uLib.Folder
		lib.Priority = uLib.Priority
		lib.PathMasks = uLib.PathMasks
		lib.Pipeline = uLib.Pipeline

		td, err := time.ParseDuration(uLib.FsCheckInterval)
		if err == nil {
			lib.FsCheckInterval = td
		}

		if err = lib.Update(); err != nil {
			logger.Error(err.Error())
			serverError(w, r, err.Error())
			return
		}

		w.Write([]byte(""))
		return
	case http.MethodDelete:
		if err = lib.Delete(); err != nil {
			logger.Error(err.Error())
			serverError(w, r, err.Error())
			return
		}
		w.Write([]byte(""))
	default:
		methodForbidden(w, r)
		return
	}

	w.Write([]byte(""))
}

func registerWebAPIv1Handlers() {
	r := newSubRouter("/api/web/v1")

	r.HandleFunc("/running", getRunning)
	r.HandleFunc("/queue", getQueue)
	r.HandleFunc("/history", getHistory)
	r.HandleFunc("/settings", settings)
	r.HandleFunc("/waitingrunners", getWaitingRunners)
	r.HandleFunc("/libraries", getAllLibraryIDs)
	r.HandleFunc("/library/", handleLibrary)
}
