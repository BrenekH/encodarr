package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/BrenekH/logange"
	"github.com/BrenekH/project-redcedar-controller/options"
)

var logger logange.Logger

func init() {
	logger = logange.NewLogger("httpserver")
}

func serverError(w http.ResponseWriter, r *http.Request, reason string) {
	logger.Warn(fmt.Sprintf("Responding to an HTTP request with code 500 because: %v", reason))
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
	wg.Add(1) // This is done in the function rather than outside so that we can easily comment out this function in main.go
	defer wg.Done()

	registerWebInterfaceHandlers()
	registerWebAPIv1Handlers()
	registerRunnerAPIv1Handlers()

	logger.Debug("Server starting")

	httpServerExitDone := &sync.WaitGroup{}

	httpServerExitDone.Add(1)
	srv := startHTTPServer(httpServerExitDone)

	logger.Info("Server started")

	<-*stopChan

	logger.Debug("Stopping HTTP server")

	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer ctxCancel()
	if err := srv.Shutdown(ctx); err != nil {
		panic(err) // Failure/timeout shutting down the server gracefully
	}

	httpServerExitDone.Wait()

	logger.Info("Server fully stopped")
}

func startHTTPServer(wg *sync.WaitGroup) *http.Server {
	srv := &http.Server{Addr: fmt.Sprintf(":%v", options.Port())}

	go func() {
		defer wg.Done()

		// Always returns error. ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// Unexpected error. port in use?
			logger.Critical(fmt.Sprintf("ListenAndServe(): %v", err))
		}
	}()

	// Returning reference so caller can call Shutdown()
	return srv
}
