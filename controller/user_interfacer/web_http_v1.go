package user_interfacer

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/BrenekH/encodarr/controller"
)

//go:embed webfiles
var webfiles embed.FS

func NewWebHTTPv1(logger controller.Logger, httpServer controller.HTTPServer, ss controller.SettingsStorer, ds controller.UserInterfacerDataStorer, useOsFs bool) WebHTTPv1 {
	return WebHTTPv1{
		logger:     logger,
		httpServer: httpServer,
		useOsFs:    useOsFs,
		ss:         ss,
		ds:         ds,
	}
}

type WebHTTPv1 struct {
	logger     controller.Logger
	httpServer controller.HTTPServer
	useOsFs    bool
	ss         controller.SettingsStorer
	ds         controller.UserInterfacerDataStorer

	waitingRunnersCache []string
}

func (w *WebHTTPv1) Start(ctx *context.Context, wg *sync.WaitGroup) {
	w.httpServer.Start(ctx, wg)

	fSys, err := fs.Sub(webfiles, "webfiles")
	if err != nil {
		w.logger.Critical(err.Error())
	}

	w.httpServer.Handle("/", http.FileServer(http.FS(fSys)))

	// React app handlers
	w.httpServer.HandleFunc("/running", w.nonRootIndexHandler)
	w.httpServer.HandleFunc("/libraries", w.nonRootIndexHandler)
	w.httpServer.HandleFunc("/history", w.nonRootIndexHandler)
	w.httpServer.HandleFunc("/settings", w.nonRootIndexHandler)

	// API Handlers
	w.httpServer.HandleFunc("/api/v1/web/running", w.getRunning)
	w.httpServer.HandleFunc("/api/v1/web/history", w.getHistory)
	w.httpServer.HandleFunc("/api/v1/web/settings", w.settings)
	w.httpServer.HandleFunc("/api/v1/web/waitingrunners", w.getWaitingRunners)
	w.httpServer.HandleFunc("/api/v1/web/libraries", w.getAllLibraryIDs)
	w.httpServer.HandleFunc("/api/v1/web/library/", w.handleLibrary)
}

func (w *WebHTTPv1) NewLibrarySettings() (m map[int]controller.Library) {
	w.logger.Critical("Not implemented")
	// TODO: Implement
	return
}

func (w *WebHTTPv1) SetLibrarySettings([]controller.Library) {
	w.logger.Critical("Not implemented")
	// TODO: Implement
}

func (w *WebHTTPv1) SetWaitingRunners(runnerNames []string) {
	w.waitingRunnersCache = runnerNames
}

// nonRootIndexHandler serves up the index files for /running, /libraries, /history, and /settings.
func (a *WebHTTPv1) nonRootIndexHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		indexFileData, err := webfiles.ReadFile("webfiles/index.html")
		if err != nil {
			a.logger.Error("Could not read 'webfiles/index.html' because of error: %v", err)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write(indexFileData)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// getRunning is a HTTP handler that returns the current running jobs in a JSON response.
func (a *WebHTTPv1) getRunning(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		dJobs, err := a.ds.DispatchedJobs()
		rJSONResp := runningJSONResponse{
			DispatchedJobs: filterDispatchedJobs(dJobs),
		}

		runningJSONBytes, err := json.Marshal(rJSONResp)
		if err != nil {
			a.logger.Error("error marshaling Job queue to json: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(runningJSONBytes)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// getHistory is a HTTP handler that returns the current history in a JSON response.
func (a *WebHTTPv1) getHistory(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		historyEntries, err := a.ds.HistoryEntries()
		if err != nil {
			a.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		h := make([]humanizedHistoryEntry, len(historyEntries))

		// Change datetime into human-readable format
		for i, v := range historyEntries {
			dt := v.DateTimeCompleted
			h[i] = humanizedHistoryEntry{
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
			a.logger.Error("error marshaling Job histroy to json: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(historyJSONBytes)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a *WebHTTPv1) settings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rS := settingsJSON{
			HealthCheckInterval: time.Duration(a.ss.HealthCheckInterval()).String(),
			HealthCheckTimeout:  time.Duration(a.ss.HealthCheckTimeout()).String(),
			LogVerbosity:        a.ss.LogVerbosity(),
		}
		b, err := json.Marshal(rS)
		if err != nil {
			a.logger.Error("failed to marshal settingsJSON: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(b)
	case http.MethodPut:
		b, err := io.ReadAll(r.Body)
		if err != nil {
			a.logger.Error(fmt.Sprintf("Failed to read request body: %v", err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		rS := settingsJSON{
			HealthCheckInterval: time.Duration(a.ss.HealthCheckInterval()).String(),
			HealthCheckTimeout:  time.Duration(a.ss.HealthCheckTimeout()).String(),
			LogVerbosity:        a.ss.LogVerbosity(),
		}
		err = json.Unmarshal(b, &rS)
		if err != nil {
			a.logger.Error(fmt.Sprintf("Failed to unmarshal settings put request body: %v", err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		td, err := time.ParseDuration(rS.HealthCheckInterval)
		if err == nil {
			a.ss.SetHealthCheckInterval(uint64(td))
		}

		td, err = time.ParseDuration(rS.HealthCheckTimeout)
		if err == nil {
			a.ss.SetHealthCheckTimeout(uint64(td))
		}

		a.ss.SetLogVerbosity(rS.LogVerbosity)

		err = a.ss.Save()
		if err != nil {
			a.logger.Error(err.Error())
		}

		w.WriteHeader(http.StatusCreated)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// getWaitingRunners is a HTTP handler that returns all runners waiting for a job in a JSON response.
func (a *WebHTTPv1) getWaitingRunners(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		b, err := json.Marshal(struct {
			Runners []string `json:"Runners"`
		}{a.waitingRunnersCache})
		if err != nil {
			a.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(b)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// getAllLibraryIDs is a HTTP handler that returns all of the library's IDs
func (a *WebHTTPv1) getAllLibraryIDs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		allLibs, err := libraries.All()
		if err != nil {
			a.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ids := make([]int, len(allLibs))
		for k, v := range allLibs {
			ids[k] = v.ID
		}

		b, err := json.Marshal(struct{ IDs []int }{ids})
		if err != nil {
			a.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(b)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// handleLibrary is a HTTP handler than takes care of the management of a Library
func (a *WebHTTPv1) handleLibrary(w http.ResponseWriter, r *http.Request) {
	libraryID := r.URL.Path[len("/api/web/v1/library/"):]

	if libraryID == "new" && r.Method == http.MethodPost {
		readBytes, err := io.ReadAll(r.Body)
		if err != nil {
			a.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		interimNewLib := libraryJSON{}
		err = json.Unmarshal(readBytes, &interimNewLib)
		if err != nil {
			a.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
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
			a.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("/api/web/v1/library/" + newLib.ID)) // TODO: Add ip/hostname to response
		return
	}

	// Transform the string libraryID into an int intLibID
	temp, err := strconv.ParseInt(libraryID, 0, 0)
	if err != nil {
		a.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	intLibID := int(temp)

	// Validate libraryID
	lib := libraries.Library{ID: intLibID}
	if err = lib.Get(); err != nil {
		a.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodGet:
		toSend := libraryJSON{lib.ID, lib.Folder, lib.Priority, lib.FsCheckInterval.String(), lib.Pipeline, lib.Queue, lib.PathMasks}
		b, err := json.Marshal(toSend)
		if err != nil {
			a.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
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
			a.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		uLib := libraryJSON{}
		err = json.Unmarshal(readBytes, &uLib)
		if err != nil {
			a.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
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
			a.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	case http.MethodDelete:
		if err = lib.Delete(); err != nil {
			a.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
