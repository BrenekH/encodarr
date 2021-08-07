package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/BrenekH/encodarr/controller"
	"github.com/BrenekH/encodarr/controller/globals"
)

// NewServer returns a new Server.
func NewServer(logger controller.Logger, port string, webAPIVersions, runnerAPIVersions []string) Server {
	registerVersionHandlers(webAPIVersions, runnerAPIVersions)

	return Server{
		port:   port,
		logger: logger,
	}
}

// Server allows multiple locations to start the same HTTP server.
type Server struct {
	serverAlreadyStarted bool
	port                 string
	logger               controller.Logger

	srv *http.Server
}

// Start starts the http server which will exit when ctx is closed. Calling Start more than once results in a no-op.
// The passed sync.WaitGroup should not have the Add method called before passing to Start.
func (s *Server) Start(ctx *context.Context, wg *sync.WaitGroup) {
	if s.serverAlreadyStarted {
		return
	}
	s.serverAlreadyStarted = true

	httpServerExitDone := sync.WaitGroup{}

	httpServerExitDone.Add(1)
	s.srv = startListenAndServer(wg, s.logger, s.port)

	wg.Add(1)
	go func() {
		defer wg.Done()

		<-(*ctx).Done()

		shutdownCtx, ctxCancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
		defer ctxCancel()
		if err := s.srv.Shutdown(shutdownCtx); err != nil {
			s.logger.Critical("%v", err) // Failure/timeout shutting down the server gracefully
		}

		httpServerExitDone.Wait()
	}()
}

// Handle wraps net/http.Handle.
func (s *Server) Handle(pattern string, handler http.Handler) {
	http.Handle(pattern, handler)
}

// HandleFunc wraps net/http.HandleFunc.
func (s *Server) HandleFunc(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, handlerFunc)
}

func startListenAndServer(wg *sync.WaitGroup, logger controller.Logger, port string) *http.Server {
	srv := &http.Server{Addr: fmt.Sprintf(":%v", port)}

	go func() {
		defer wg.Done()

		// Always returns error. ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// Unexpected error. port in use?
			logger.Error("unexpected error: %v\n", err)
		}
	}()

	// Returning reference so caller can call Shutdown()
	return srv
}

func registerVersionHandlers(webVersions, runnerVersions []string) {
	// Controller version
	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(globals.Version))
	})

	// Web and Runner versions
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		respStruct := struct {
			Web struct {
				Versions []string `json:"versions"`
			} `json:"web"`
			Runner struct {
				Versions []string `json:"versions"`
			} `json:"runner"`
		}{
			Web: struct {
				Versions []string `json:"versions"`
			}{webVersions},
			Runner: struct {
				Versions []string `json:"versions"`
			}{runnerVersions},
		}

		b, _ := json.Marshal(respStruct)
		w.Write(b)
	})

	// Web versions
	http.HandleFunc("/api/web", func(w http.ResponseWriter, r *http.Request) {
		respStruct := struct {
			Versions []string `json:"versions"`
		}{
			Versions: webVersions,
		}

		b, _ := json.Marshal(respStruct)
		w.Write(b)
	})

	// Runner versions
	http.HandleFunc("/api/runner", func(w http.ResponseWriter, r *http.Request) {
		respStruct := struct {
			Versions []string `json:"versions"`
		}{
			Versions: webVersions,
		}

		b, _ := json.Marshal(respStruct)
		w.Write(b)
	})
}
