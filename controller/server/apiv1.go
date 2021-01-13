package server

import (
	"net/http"
)

// Web interface API handlers
// TODO: Get running jobs
// TODO: Get queue
// TODO: Get history
func apiSample(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><head><title>API Test - Project RedCedar</title></head><body><h4>Hello, World!</h4></body></html>`))
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

	r.HandleFunc("/sample", apiSample)
}
