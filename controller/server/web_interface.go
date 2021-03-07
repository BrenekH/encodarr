package server

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed webfiles
var webfiles embed.FS

func nonRootIndexHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		indexFileData, err := webfiles.ReadFile("webfiles/index.html")
		if err != nil {
			serverError(w, r, fmt.Sprintf("Could not read 'webfiles/index.html' because of error: %v", err))
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write(indexFileData)
	default:
		methodForbidden(w, r)
	}
}

func registerWebInterfaceHandlers() {
	fSys, err := fs.Sub(webfiles, "webfiles")
	if err != nil {
		panic(err)
	}
	http.Handle("/", http.FileServer(http.FS(fSys)))

	// Non-root handlers (/running, /queue, /history, and /settings should all send index.html, but by default they don't)
	http.HandleFunc("/running", nonRootIndexHandler)
	http.HandleFunc("/queue", nonRootIndexHandler)
	http.HandleFunc("/history", nonRootIndexHandler)
	http.HandleFunc("/settings", nonRootIndexHandler)
}
