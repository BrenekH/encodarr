package userinterfacer

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

// NewWebHTTPv1 uses the provided arguments to instantiate a new WebHTTPv1 struct and return it.
func NewWebHTTPv1(logger controller.Logger, httpServer controller.HTTPServer, ss controller.SettingsStorer, ds controller.UserInterfacerDataStorer, useOsFs bool) WebHTTPv1 {
	return WebHTTPv1{
		logger:     logger,
		httpServer: httpServer,
		useOsFs:    useOsFs,
		ss:         ss,
		ds:         ds,

		waitingRunnersCache: make([]string, 0),
		libraryCache:        []controller.Library{},
		libSettingsUpdates:  map[int]controller.Library{},
	}
}

// WebHTTPv1 satisfies the controller.UserInterfacer interface using http.
type WebHTTPv1 struct {
	logger     controller.Logger
	httpServer controller.HTTPServer
	useOsFs    bool
	ss         controller.SettingsStorer
	ds         controller.UserInterfacerDataStorer

	waitingRunnersCache []string
	libraryCache        []controller.Library
	libSettingsUpdates  map[int]controller.Library
}

// Start starts the http server without blocking the thread.
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
	w.httpServer.HandleFunc("/api/web/v1/running", w.getRunning)
	w.httpServer.HandleFunc("/api/web/v1/history", w.getHistory)
	w.httpServer.HandleFunc("/api/web/v1/settings", w.settings)
	w.httpServer.HandleFunc("/api/web/v1/waitingrunners", w.getWaitingRunners)
	w.httpServer.HandleFunc("/api/web/v1/libraries", w.getAllLibraryIDs)
	w.httpServer.HandleFunc("/api/web/v1/library/", w.handleLibrary)
}

// NewLibrarySettings returns a new library settings the user may have set.
func (w *WebHTTPv1) NewLibrarySettings() map[int]controller.Library {
	copy := w.libSettingsUpdates
	w.libSettingsUpdates = map[int]controller.Library{}
	return copy
}

// SetLibrarySettings sets the library settings cache to be shown to the user.
func (w *WebHTTPv1) SetLibrarySettings(libs []controller.Library) {
	w.libraryCache = libs
}

// SetWaitingRunners sets the list of waiting runners in memory so that it can be shown to the user.
func (w *WebHTTPv1) SetWaitingRunners(runnerNames []string) {
	// We have to make and copy runnerNames here so that when we marshal the w.waitingRunnersCache slice to json, it isn't null.
	temp := make([]string, len(runnerNames))
	copy(temp, runnerNames)
	w.waitingRunnersCache = temp
}

// nonRootIndexHandler serves up the index files for /running, /libraries, /history, and /settings.
func (w *WebHTTPv1) nonRootIndexHandler(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		indexFileData, err := webfiles.ReadFile("webfiles/index.html")
		if err != nil {
			w.logger.Error("Could not read 'webfiles/index.html' because of error: %v", err)
			return
		}
		rw.Header().Set("Content-Type", "text/html")
		rw.WriteHeader(http.StatusOK)
		rw.Write(indexFileData)
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// getRunning is a HTTP handler that returns the current running jobs in a JSON response.
func (w *WebHTTPv1) getRunning(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		dJobs, err := w.ds.DispatchedJobs()
		if err != nil {
			w.logger.Error(err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rJSONResp := runningJSONResponse{
			DispatchedJobs: filterDispatchedJobs(dJobs),
		}

		runningJSONBytes, err := json.Marshal(rJSONResp)
		if err != nil {
			w.logger.Error("error marshaling Job queue to json: %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.Write(runningJSONBytes)
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// getHistory is a HTTP handler that returns the current history in a JSON response.
func (w *WebHTTPv1) getHistory(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		historyEntries, err := w.ds.HistoryEntries()
		if err != nil {
			w.logger.Error(err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
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
			w.logger.Error("error marshaling Job histroy to json: %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.Write(historyJSONBytes)
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// settings is a HTTP handler for both setting and getting the current Controller settings.
func (w *WebHTTPv1) settings(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rS := settingsJSON{
			HealthCheckInterval: time.Duration(w.ss.HealthCheckInterval()).String(),
			HealthCheckTimeout:  time.Duration(w.ss.HealthCheckTimeout()).String(),
			LogVerbosity:        w.ss.LogVerbosity(),
		}
		b, err := json.Marshal(rS)
		if err != nil {
			w.logger.Error("failed to marshal settingsJSON: %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Write(b)
	case http.MethodPut:
		b, err := io.ReadAll(r.Body)
		if err != nil {
			w.logger.Error(fmt.Sprintf("Failed to read request body: %v", err))
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		rS := settingsJSON{
			HealthCheckInterval: time.Duration(w.ss.HealthCheckInterval()).String(),
			HealthCheckTimeout:  time.Duration(w.ss.HealthCheckTimeout()).String(),
			LogVerbosity:        w.ss.LogVerbosity(),
		}
		err = json.Unmarshal(b, &rS)
		if err != nil {
			w.logger.Error(fmt.Sprintf("Failed to unmarshal settings put request body: %v", err))
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		td, err := time.ParseDuration(rS.HealthCheckInterval)
		if err == nil {
			w.ss.SetHealthCheckInterval(uint64(td))
		}

		td, err = time.ParseDuration(rS.HealthCheckTimeout)
		if err == nil {
			w.ss.SetHealthCheckTimeout(uint64(td))
		}

		w.ss.SetLogVerbosity(rS.LogVerbosity)

		err = w.ss.Save()
		if err != nil {
			w.logger.Error(err.Error())
		}

		rw.WriteHeader(http.StatusCreated)
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// getWaitingRunners is a HTTP handler that returns all runners waiting for a job in a JSON response.
func (w *WebHTTPv1) getWaitingRunners(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		b, err := json.Marshal(struct {
			Runners []string `json:"Runners"`
		}{w.waitingRunnersCache})
		if err != nil {
			w.logger.Error(err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Write(b)
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// getAllLibraryIDs is a HTTP handler that returns all of the library's IDs
func (w *WebHTTPv1) getAllLibraryIDs(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ids := make([]int, len(w.libraryCache))
		for k, v := range w.libraryCache {
			ids[k] = v.ID
		}

		b, err := json.Marshal(struct{ IDs []int }{ids})
		if err != nil {
			w.logger.Error(err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Write(b)
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// handleLibrary is a HTTP handler than takes care of the management of a Library
func (w *WebHTTPv1) handleLibrary(rw http.ResponseWriter, r *http.Request) {
	libraryID := r.URL.Path[len("/api/web/v1/library/"):]

	if libraryID == "new" && r.Method == http.MethodPost {
		readBytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.logger.Error(err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		interimNewLib := interimLibraryJSON{}
		err = json.Unmarshal(readBytes, &interimNewLib)
		if err != nil {
			w.logger.Error(err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		newLib := controller.Library{
			Folder:    interimNewLib.Folder,
			Priority:  interimNewLib.Priority,
			PathMasks: interimNewLib.PathMasks,
		}

		td, err := time.ParseDuration(interimNewLib.FsCheckInterval)
		if err == nil {
			newLib.FsCheckInterval = td
		}

		// Create map of library IDs (for fast valid ID lookup)
		libIDMap := map[int]struct{}{}
		for _, v := range w.libraryCache {
			libIDMap[v.ID] = struct{}{}
		}

		// Find valid ID
		var validID int
		for i := 0; i < 10_000; i++ {
			_, ok := libIDMap[i]
			if !ok {
				validID = i
				break
			}
		}
		newLib.ID = validID

		// Put new lib in a.libSettingsUpdates with located valid id
		w.libSettingsUpdates[newLib.ID] = newLib

		rw.WriteHeader(http.StatusCreated)
		rw.Write([]byte(fmt.Sprintf("http://%v/api/web/v1/library/%v", r.Host, newLib.ID)))
		return
	}

	// Transform the string libraryID into an int intLibID
	temp, err := strconv.ParseInt(libraryID, 0, 0)
	if err != nil {
		w.logger.Error(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	intLibID := int(temp)

	// Validate libraryID (and pull the matching library out of the cache).
	var lib controller.Library
	validID := false
	for _, v := range w.libraryCache {
		if intLibID == v.ID {
			validID = true
			lib = v
			break
		}
	}

	if !validID {
		w.logger.Error("invalid library ID '%v' requested", intLibID)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		toSend := interimLibraryJSON{lib.ID, lib.Folder, lib.Priority, lib.FsCheckInterval.String(), lib.Queue, lib.PathMasks, lib.CommandDeciderSettings}
		b, err := json.Marshal(toSend)
		if err != nil {
			w.logger.Error(err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.Write(b)
	case http.MethodPut:
		// Technically, there is a security flaw where an attacker can set the id in their request
		// to a different library and overwrite a different library, but it's not like this API is locked down at all
		// so does it really matter?
		readBytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.logger.Error(err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		uLib := interimLibraryJSON{}
		err = json.Unmarshal(readBytes, &uLib)
		if err != nil {
			w.logger.Error(err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		lib.Folder = uLib.Folder
		lib.Priority = uLib.Priority
		lib.PathMasks = uLib.PathMasks
		lib.CommandDeciderSettings = uLib.CommandDeciderSettings

		td, err := time.ParseDuration(uLib.FsCheckInterval)
		if err == nil {
			lib.FsCheckInterval = td
		}

		// Add lib to response of UI.NewLibrarySettings
		w.libSettingsUpdates[lib.ID] = lib

		rw.WriteHeader(http.StatusNoContent)
	case http.MethodDelete:
		if err = w.ds.DeleteLibrary(lib.ID); err != nil {
			w.logger.Error(err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusNoContent)
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}
