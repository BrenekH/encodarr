package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/BrenekH/project-redcedar-controller/controller"
)

type incomingJobStatus struct {
	UUID   string               `json:"uuid"`
	Status controller.JobStatus `json:"status"`
}

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
			serverError(w, r, fmt.Sprintf("Runner API v1: Error marshaling Job to json: %v", err))
		}
		w.Header().Set("x-rc-job-info", string(jobJSONBytes))

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
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Runner API v1: Error reading job status body: %v", err)
		}

		ijs := incomingJobStatus{}
		err = json.Unmarshal(b, &ijs)
		if err != nil {
			log.Printf("Runner API v1: Error unmarshalling into struct: %v", err)
		}

		err = controller.DispatchedJobs.UpdateStatus(ijs.UUID, ijs.Status)
		if err != nil { // Since I wrote UpdateStatus, I know that if it errors at all, it's an issue with the UUID
			//! Technically this is the clients fault, not the server's, so a different HTTP code should be sent
			serverError(w, r, fmt.Sprintf("Runner API v1: Error updating status of job with uuid '%v': %v", ijs.UUID, err))
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	default:
		methodForbidden(w, r)
	}
}

// TODO: Complete post job complete
func postJobComplete(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h := r.Header.Get("x-rc-history-entry")
		if h == "" {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid header 'x-rc-history-entry'"))
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
			_ = fileHeader // File header could be useful for naming later
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

		controller.CompletedRequestChannel <- jcr

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
