package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/BrenekH/project-redcedar-controller/controller"
)

func getNewJob(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		requestChannel := make(chan controller.Job, 1)
		controller.JobRequestChannel <- controller.JobRequest{RunnerName: r.Header.Get("redcedar-runner-name"), ReturnChannel: &requestChannel}
		jobToSend, ok := <-requestChannel

		if ok == false {
			serverError(w, r, "Server shutdown")
			return
		}

		w.Header().Set("Content-Type", inferMIMETypeFromExt(filepath.Ext(jobToSend.Path)))

		// Marshal Job into json to be sent in a header
		jobJSONBytes, err := json.Marshal(jobToSend)
		if err != nil {
			serverError(w, r, fmt.Sprintf("APIv1: Error marshaling Job to json: %v", err))
		}
		w.Header().Set("x-rc-job-info", string(jobJSONBytes))

		// Send file to Runner
		file, err := os.Open(jobToSend.Path)
		if err != nil {
			serverError(w, r, fmt.Sprintf("APIv1: Error opening %v: %v", jobToSend.Path, err))
			return
		}
		defer file.Close()

		bufferSize := 1024
		buffer := make([]byte, bufferSize)

		for {
			bytesRead, err := file.Read(buffer)
			if err != nil {
				if err != io.EOF {
					serverError(w, r, fmt.Sprintf("APIv1: Error writing to HTTP Writer: %v", err))
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

// TODO: Complete post job status
func postJobStatus(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"test": true}`))
	default:
		methodForbidden(w, r)
	}
}

// TODO: Complete post job complete
func postJobComplete(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(""))
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
