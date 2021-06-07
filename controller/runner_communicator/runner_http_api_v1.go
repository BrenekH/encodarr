package runner_communicator

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/BrenekH/encodarr/controller"
)

func NewRunnerHTTPApiV1(logger controller.Logger, httpServer controller.HTTPServer, ds controller.RunnerCommunicatorDataStorer) RunnerHTTPApiV1 {
	return RunnerHTTPApiV1{
		logger:         logger,
		httpServer:     httpServer,
		ds:             ds,
		nullifiedUUIDs: make([]controller.UUID, 0),
	}
}

type RunnerHTTPApiV1 struct {
	logger     controller.Logger
	httpServer controller.HTTPServer
	ds         controller.RunnerCommunicatorDataStorer

	nullifiedUUIDs []controller.UUID
}

func (r *RunnerHTTPApiV1) Start(ctx *context.Context, wg *sync.WaitGroup) {
	r.httpServer.Start(ctx, wg)

	// Add handlers to r.httpServer
	r.httpServer.HandleFunc("/api/runner/v1/job/request", r.requestJob)
	r.httpServer.HandleFunc("/api/runner/v1/job/status", r.jobStatus)
	r.httpServer.HandleFunc("/api/runner/v1/job/complete", r.jobComplete)
}

func (r *RunnerHTTPApiV1) CompletedJobs() (j []controller.CompletedJob) {
	r.logger.Critical("Not Implemented")
	// TODO: Implement

	// NOTE: Use a channel to transfer all completed job requests from the HTTP handler to this function.

	return
}

func (r *RunnerHTTPApiV1) NewJob(controller.Job) {
	r.logger.Critical("Not Implemented")
	// TODO: Implement
}

func (r *RunnerHTTPApiV1) NeedNewJob() (b bool) {
	r.logger.Critical("Not Implemented")
	// TODO: Implement
	return
}

func (r *RunnerHTTPApiV1) NullifyUUIDs(uuids []controller.UUID) {
	r.nullifiedUUIDs = append(r.nullifiedUUIDs, uuids...)
}

func (r *RunnerHTTPApiV1) WaitingRunners() (runnerNames []string) {
	r.logger.Critical("Not Implemented")
	// TODO: Implement
	return
}

func (a *RunnerHTTPApiV1) requestJob(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// TODO: Gather Runner name from HTTP headers
		// TODO: Add callback channel to waiting runners queue
		// TODO: Check for a returned job
		// TODO: Also check for a connection close using r.Context(). Remove from waiting runners if becomes done.
		// TODO: Send back header info
		// TODO: Respond with file
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a *RunnerHTTPApiV1) jobStatus(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		b, err := io.ReadAll(r.Body)
		if err != nil {
			a.logger.Error("error reading job status body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ijs := incomingJobStatus{}
		if err = json.Unmarshal(b, &ijs); err != nil {
			a.logger.Error("error unmarshalling into struct: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Check ijs.UUID against nullified UUIDs
		for _, v := range a.nullifiedUUIDs {
			if v == ijs.UUID {
				w.WriteHeader(http.StatusConflict) // Send the 409 error code to signal to the Runner that the job has been nullified.
				return
			}
		}

		// Get existing DispatchedJob from datastore
		dJob, err := a.ds.DispatchedJob(ijs.UUID)
		if err != nil {
			a.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Update DispatchedJob.Status
		dJob.Status = ijs.Status

		// Store DispatchedJob into datastore
		if err = a.ds.SaveDispatchedJob(dJob); err != nil {
			a.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a *RunnerHTTPApiV1) jobComplete(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// TODO: Implement
		// Goals: Check UUID for nullified status, save file that client is posting, add job to channel for CompletedJobs to pick up from

		//? Probably should detect a client disconnect and disregard any data it sent
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

type incomingJobStatus struct {
	UUID   controller.UUID      `json:"uuid"`
	Status controller.JobStatus `json:"status"`
}
