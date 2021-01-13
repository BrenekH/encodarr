package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
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

// RunHTTPServer runs the HTTP server for Project RedCedar.
func RunHTTPServer(stopChan *chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	registerWebInterfaceHandlers()
	registerAPIv1Handlers()

	log.Printf("HTTP Server: Server starting")

	httpServerExitDone := &sync.WaitGroup{}

	httpServerExitDone.Add(1)
	srv := startHTTPServer(httpServerExitDone)

	log.Printf("HTTP Server: Server started")

	<-*stopChan

	log.Printf("HTTP Server: Stopping HTTP server")

	// TODO: Replace TODO context with something more appropriate
	if err := srv.Shutdown(context.TODO()); err != nil {
		panic(err) // Failure/timeout shutting down the server gracefully
	}

	httpServerExitDone.Wait()

	log.Printf("HTTP Server: Server fully stopped")
}

func startHTTPServer(wg *sync.WaitGroup) *http.Server {
	srv := &http.Server{Addr: ":8080"}

	go func() {
		defer wg.Done()

		// Always returns error. ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// Unexpected error. port in use?
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// Returning reference so caller can call Shutdown()
	return srv
}
