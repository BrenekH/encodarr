package user_interfacer

import (
	"context"
	"embed"
	"io/fs"
	"net/http"
	"sync"

	"github.com/BrenekH/encodarr/controller"
)

//go:embed webfiles
var webfiles embed.FS

func NewWebHTTPApiV1(logger controller.Logger, httpServer controller.HTTPServer, useOsFs bool) WebHTTPApiV1 {
	return WebHTTPApiV1{logger: logger, httpServer: httpServer, useOsFs: useOsFs}
}

type WebHTTPApiV1 struct {
	logger     controller.Logger
	httpServer controller.HTTPServer
	useOsFs    bool
}

func (w *WebHTTPApiV1) Start(ctx *context.Context, wg *sync.WaitGroup) {
	w.httpServer.Start(ctx, wg)

	fSys, err := fs.Sub(webfiles, "webfiles")
	if err != nil {
		w.logger.Critical(err.Error())
	}

	w.httpServer.Handle("/", http.FileServer(http.FS(fSys)))

	w.httpServer.HandleFunc("/running", w.nonRootIndexHandler)
	w.httpServer.HandleFunc("/libraries", w.nonRootIndexHandler)
	w.httpServer.HandleFunc("/history", w.nonRootIndexHandler)
	w.httpServer.HandleFunc("/settings", w.nonRootIndexHandler)

	// TODO: Add API handlers to w.httpServer
}

func (w *WebHTTPApiV1) NewLibrarySettings() (m map[int]controller.Library) {
	w.logger.Critical("Not implemented")
	// TODO: Implement
	return
}

func (w *WebHTTPApiV1) SetLibrarySettings([]controller.Library) {
	w.logger.Critical("Not implemented")
	// TODO: Implement
}

func (w *WebHTTPApiV1) SetWaitingRunners(runnerNames []string) {
	w.logger.Critical("Not implemented")
	// TODO: Implement
}

// nonRootIndexHandler serves up the index files for /running, /libraries, /history, and /settings.
func (a *WebHTTPApiV1) nonRootIndexHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		indexFileData, err := webfiles.ReadFile("webfiles/index.html")
		if err != nil {
			a.logger.Error("Could not read 'webfiles/index.html' because of error: %v", err)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write(indexFileData)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
