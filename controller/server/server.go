package server

import (
	"fmt"
	"log"
	"net/http"
)

func serverError(w http.ResponseWriter, r *http.Request, reason string) {
	fmt.Println(reason)
	w.WriteHeader(500)
	w.Write([]byte(`<html><body>Code 500: Server Error</body></html>`))
}

// RunHTTPServer runs the HTTP server for Project RedCedar in a blocking manner.
func RunHTTPServer() {
	envPort := "8080"
	registerWebInterfaceHandlers()
	registerAPIv1Handlers()
	log.Fatal(http.ListenAndServe(":"+envPort, nil))
}
