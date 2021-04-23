package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/BrenekH/encodarr/controller/controller"
	"github.com/BrenekH/encodarr/controller/db/dispatched"
)

type incomingJobStatus struct {
	UUID   string               `json:"uuid"`
	Status dispatched.JobStatus `json:"status"`
}

func getNewJob(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		logger.Info(fmt.Sprintf("Received new job request from %v @ %v", r.Header.Get("X-Encodarr-Runner-Name"), r.RemoteAddr))
		requestChannel := make(chan dispatched.Job, 1)
		controller.JobRequestChannel <- controller.JobRequest{RunnerName: r.Header.Get("X-Encodarr-Runner-Name"), ReturnChannel: &requestChannel}
		jobToSend, ok := <-requestChannel

		if !ok {
			serverError(w, r, "Server shutdown")
			return
		}

		w.Header().Set("Content-Type", inferMIMETypeFromExt(filepath.Ext(jobToSend.Path)))

		// Marshal Job into json to be sent in a header
		jobJSONBytes, err := json.Marshal(jobToSend)
		if err != nil {
			serverError(w, r, fmt.Sprintf("Runner API v1: Error marshaling Job to json: %v", err))
		}
		w.Header().Set("X-Encodarr-Job-Info", string(jobJSONBytes))

		// Send file to Runner
		file, err := os.Open(jobToSend.Path)
		if err != nil {
			serverError(w, r, fmt.Sprintf("Runner API v1: Error opening %v: %v", jobToSend.Path, err))
			return
		}
		defer file.Close()

		bufferSize := 1024
		buffer := make([]byte, bufferSize)

		for {
			bytesRead, err := file.Read(buffer)
			if err != nil {
				if err != io.EOF {
					serverError(w, r, fmt.Sprintf("Runner API v1: Error writing to HTTP Writer: %v", err))
					return
				}
				break
			}
			w.Write(buffer[:bytesRead])
		}

	default:
		methodForbidden(w, r)
	}
}

func postJobStatus(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		b, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error(fmt.Sprintf("Runner API v1: Error reading job status body: %v", err))
		}

		ijs := incomingJobStatus{}
		err = json.Unmarshal(b, &ijs)
		if err != nil {
			logger.Error(fmt.Sprintf("Runner API v1: Error unmarshalling into struct: %v", err))
		}

		dJob := dispatched.DJob{UUID: ijs.UUID}
		err = dJob.Get()
		if err != nil {
			logger.Error(err.Error())
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusConflict) // Sends the 409 code to signal runners to abandon the job
			w.Write([]byte(""))
			return
		}

		dJob.Status = ijs.Status
		dJob.LastUpdated = time.Now()
		err = dJob.Update()
		if err != nil {
			logger.Error(fmt.Sprintf("Error saving dispatched jobs: %v", err.Error()))
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	default:
		methodForbidden(w, r)
	}
}

func postJobComplete(w http.ResponseWriter, r *http.Request) {
	logger.Debug(fmt.Sprintf("Received /api/runner/v1/job/complete from %v", r.RemoteAddr))
	switch r.Method {
	case http.MethodPost:
		h := r.Header.Get("X-Encodarr-History-Entry")
		if h == "" {
			logger.Debug("Received invalid history entry")
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid header 'X-Encodarr-History-Entry'"))
			return
		}

		jcr := controller.JobCompleteRequest{}
		err := json.Unmarshal([]byte(h), &jcr)
		if err != nil {
			serverError(w, r, fmt.Sprintf("Error unmarshalling history entry: %v", err))
			return
		}

		if !jcr.Failed {
			fileReader, fileHeader, err := r.FormFile("file")
			_ = fileHeader //! File header could be useful for naming later
			if err != nil {
				serverError(w, r, fmt.Sprintf("Error accessing form file: %v", err))
				return
			}
			defer fileReader.Close()

			// Copy to intermediate file
			jcr.InFile = fmt.Sprintf("%v.import%v", jcr.UUID, path.Ext(fileHeader.Filename))
			f, err := os.Create(jcr.InFile)
			if err != nil {
				serverError(w, r, fmt.Sprintf("Error opening receiving file: %v", err))
				return
			}
			io.Copy(f, fileReader)
			f.Close()
		}

		controller.JobCompleteRequestChan <- jcr

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	default:
		methodForbidden(w, r)
	}
}

func registerRunnerAPIv1Handlers() {
	r := newSubRouter("/api/runner/v1")

	r.HandleFunc("/job/request", getNewJob)
	r.HandleFunc("/job/status", postJobStatus)
	r.HandleFunc("/job/complete", postJobComplete)
}
