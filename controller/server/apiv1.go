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

// Web interface API handlers
// TODO: Get running jobs

// TODO: Complete GET queue
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

// TODO: Complete GET history
func getHistory(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"test": true}`))
	default:
		methodForbidden(w, r)
	}
}

// Runner API handlers
// TODO: Get new job (job request)
// TODO: Post job status
// TODO: Post job complete

func registerAPIv1Handlers() {
	r := newSubRouter("/api/v1")

	r.HandleFunc("/queue", getQueue)
	r.HandleFunc("/history", getHistory)
}
