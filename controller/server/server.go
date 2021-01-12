package server

import (
	"fmt"
	"log"
	"net/http"
)

func serverError(w http.ResponseWriter, r *http.Request, reason string) {
	fmt.Println(reason)
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(500)
	w.Write([]byte(`<html><head><title>Server Error - Project RedCedar</title></head><body>Code 500: Server Error</body></html>`))
}

func methodForbidden(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(`<html><head><title>Method Forbidden - Project RedCedar</title></head><body>Code 403: Method Forbidden</body></html>`))
}

// RunHTTPServer runs the HTTP server for Project RedCedar in a blocking manner.
func RunHTTPServer() {
	envPort := "8080"
	registerWebInterfaceHandlers()
	registerAPIv1Handlers()
	log.Fatal(http.ListenAndServe(":"+envPort, nil))
}
