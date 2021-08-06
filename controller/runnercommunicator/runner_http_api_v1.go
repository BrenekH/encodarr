package runnercommunicator

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/BrenekH/encodarr/controller"
	"github.com/google/uuid"
)

// NewRunnerHTTPApiV1 returns a new RunnerHTTPApiV1.
func NewRunnerHTTPApiV1(logger controller.Logger, httpServer controller.HTTPServer, ds controller.RunnerCommunicatorDataStorer) RunnerHTTPApiV1 {
	return RunnerHTTPApiV1{
		logger:         logger,
		httpServer:     httpServer,
		ds:             ds,
		nullifiedUUIDs: make([]controller.UUID, 0),
		wrQueue:        newQueue(),
		completedJobs:  make(chan controller.CompletedJob),
	}
}

// RunnerHTTPApiV1 satisfies the controller.RunnerCommunicator interface using an HTTP server.
type RunnerHTTPApiV1 struct {
	logger     controller.Logger
	httpServer controller.HTTPServer
	ds         controller.RunnerCommunicatorDataStorer

	nullifiedUUIDs []controller.UUID
	completedJobs  chan controller.CompletedJob
	wrQueue        queue
}

// Start starts the HTTP server. It does not block the thread.
func (r *RunnerHTTPApiV1) Start(ctx *context.Context, wg *sync.WaitGroup) {
	r.httpServer.Start(ctx, wg)

	// Add handlers to r.httpServer
	r.httpServer.HandleFunc("/api/runner/v1/job/request", r.requestJob)
	r.httpServer.HandleFunc("/api/runner/v1/job/status", r.jobStatus)
	r.httpServer.HandleFunc("/api/runner/v1/job/complete", r.jobComplete)
}

// CompletedJobs returns any completed jobs as told by the Runners.
func (r *RunnerHTTPApiV1) CompletedJobs() (j []controller.CompletedJob) {
	for {
		select {
		case cJob := <-r.completedJobs:
			j = append(j, cJob)
		default:
			return // Have to use return, since break will only cancel the select clause and not the for loop.
		}
	}
}

// NewJob sends a new job to the next running in the queue.
func (r *RunnerHTTPApiV1) NewJob(cJob controller.Job) {
	wr, err := r.wrQueue.Pop()
	if err != nil {
		r.logger.Error("NewJob was called but got error from Pop: %v", err)
	}

	// Add job to dispatched jobs
	dJob := controller.DispatchedJob{
		UUID:        cJob.UUID,
		Runner:      wr.Name,
		Job:         cJob,
		Status:      controller.JobStatus{},
		LastUpdated: time.Now(),
	}
	err = r.ds.SaveDispatchedJob(dJob)
	if err != nil {
		r.logger.Error("error saving new dispatched job: %v", err)
	}

	wr.CallbackChan <- cJob
}

// NeedNewJob returns a boolean indicating whether or not there are waiting runners.
func (r *RunnerHTTPApiV1) NeedNewJob() bool {
	return !r.wrQueue.Empty()
}

// NullifyUUIDs marks a job uuid as null.
func (r *RunnerHTTPApiV1) NullifyUUIDs(uuids []controller.UUID) {
	r.nullifiedUUIDs = append(r.nullifiedUUIDs, uuids...)
}

// WaitingRunners returns a list of Runner names which are waiting for jobs.
func (r *RunnerHTTPApiV1) WaitingRunners() (runnerNames []string) {
	wr := r.wrQueue.Dequeue()

	for _, v := range wr {
		runnerNames = append(runnerNames, v.Name)
	}

	return
}

// requestJob is an HTTP handler which handles Runner job requests.
func (r *RunnerHTTPApiV1) requestJob(w http.ResponseWriter, hr *http.Request) {
	r.logger.Trace("%+v", hr)
	switch hr.Method {
	case http.MethodGet:
		// Gather Runner name from HTTP headers
		runnerName := hr.Header.Get("X-Encodarr-Runner-Name")
		r.logger.Info("Received request from %v @ %v", runnerName, hr.RemoteAddr)

		// Add callback channel to waiting runners queue
		receiveChan := make(chan controller.Job)
		requestUUID := uuid.NewString()
		r.wrQueue.Push(waitingRunner{Name: runnerName, CallbackChan: receiveChan, UUID: requestUUID})

		// Check for a returned job
		var jobToSend controller.Job
		var ok bool
		select {
		case jobToSend, ok = <-receiveChan:
			break
		case <-hr.Context().Done():
			r.wrQueue.Remove(requestUUID)
			w.WriteHeader(http.StatusGone)
			return
		}

		close(receiveChan)

		if !ok { // Server shutdown
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Marshal Job into json to be sent in a header
		jobJSONBytes, err := json.Marshal(jobToSend)
		if err != nil {
			r.logger.Error("error marshaling Job to json: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("X-Encodarr-Job-Info", string(jobJSONBytes))

		// Respond with file
		w.Header().Set("Content-Type", inferMIMETypeFromExt(filepath.Ext(jobToSend.Path)))
		file, err := os.Open(jobToSend.Path)
		if err != nil {
			r.logger.Error("error opening %v: %v", jobToSend.Path, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer file.Close()

		bufferSize := 1024
		buffer := make([]byte, bufferSize)

		for {
			bytesRead, err := file.Read(buffer)
			if err != nil {
				if err != io.EOF {
					r.logger.Error("error writing to HTTP Writer: %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				break
			}
			w.Write(buffer[:bytesRead])
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// jobStatus in an HTTP handler for status updates from Runners.
func (r *RunnerHTTPApiV1) jobStatus(w http.ResponseWriter, hr *http.Request) {
	r.logger.Trace("%+v", hr)
	switch hr.Method {
	case http.MethodPost:
		b, err := io.ReadAll(hr.Body)
		if err != nil {
			r.logger.Error("error reading job status body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ijs := incomingJobStatus{}
		if err = json.Unmarshal(b, &ijs); err != nil {
			r.logger.Error("error unmarshalling into struct: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Check ijs.UUID against nullified UUIDs
		for _, v := range r.nullifiedUUIDs {
			if v == ijs.UUID {
				w.WriteHeader(http.StatusConflict) // Send the 409 error code to signal to the Runner to indicate that the job has been nullified.
				return
			}
		}

		// Get existing DispatchedJob from datastore
		dJob, err := r.ds.DispatchedJob(ijs.UUID)
		if err != nil {
			r.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		dJob.Status = ijs.Status

		// Update the LastUpdated time so that the health check won't null this Runner
		dJob.LastUpdated = time.Now()

		// Store DispatchedJob into datastore
		if err = r.ds.SaveDispatchedJob(dJob); err != nil {
			r.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// jobComplete is an HTTP handler for Runners to upload completed jobs to.
func (r *RunnerHTTPApiV1) jobComplete(w http.ResponseWriter, hr *http.Request) {
	r.logger.Trace("%+v", hr)
	switch hr.Method {
	case http.MethodPost:
		//? Probably should detect a client disconnect and disregard any data it sent (use r.Context())

		// Read history entry from headers
		h := hr.Header.Get("X-Encodarr-History-Entry")
		if h == "" {
			r.logger.Debug("received invalid history entry")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Unmarshal history entry into usable struct
		cJob := controller.CompletedJob{}
		err := json.Unmarshal([]byte(h), &cJob)
		if err != nil {
			r.logger.Debug("error unmarshalling history entry: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// If UUID was nullified, respond with 409 error code and exit
		for _, v := range r.nullifiedUUIDs {
			if v == cJob.UUID {
				w.WriteHeader(http.StatusConflict)
				return
			}
		}

		// If job didn't fail, write file to disk
		if !cJob.Failed {
			fileReader, fileHeader, err := hr.FormFile("file")
			if err != nil {
				r.logger.Debug("error accessing form file: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			defer fileReader.Close()

			// Copy to intermediate file
			cJob.InFile = fmt.Sprintf("%v.import%v", cJob.UUID, filepath.Ext(fileHeader.Filename))
			f, err := os.Create(cJob.InFile) // TODO: Mock out for testing
			if err != nil {
				r.logger.Debug("error opening receiving file: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			io.Copy(f, fileReader)
			f.Close()
		}

		// Add controller.CompletedJob to channel for CompletedJobs to pick up from
		r.completedJobs <- cJob

		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// incomingJobStatus defines the structure of the job status from Runners.
type incomingJobStatus struct {
	UUID   controller.UUID      `json:"uuid"`
	Status controller.JobStatus `json:"status"`
}
