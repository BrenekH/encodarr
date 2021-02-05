package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// favicon is a HTTP handler for the favicon.ico file
func favicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/x-icon")

	icoData, err := ioutil.ReadFile("webfiles/favicon/favicon.ico")
	if err != nil {
		serverError(w, r, fmt.Sprintf("Could not read %v because of error: %v", r.URL, err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(icoData)
}

// resources is a HTTP handler for resource(css, js, svg) requests to the server.
func resources(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		cssRequest := strings.Contains(r.URL.String(), "css")
		jsRequest := strings.Contains(r.URL.String(), "js")
		jpgRequest := strings.Contains(r.URL.String(), "jpg") || strings.Contains(r.URL.String(), "jpeg")
		svgRequest := strings.Contains(r.URL.String(), "svg")

		if cssRequest {
			w.Header().Set("Content-Type", "text/css")
		} else if jsRequest {
			w.Header().Set("Content-Type", "text/javascript")
		} else if jpgRequest {
			w.Header().Set("Content-Type", "image/jpeg")
		} else if svgRequest {
			w.Header().Set("Content-Type", "image/svg+xml")
		} else {
			logger.Warn(fmt.Sprintf("Could not identify MIME type for resources request: %v\n", r.URL.String()))
			w.Header().Set("Content-Type", "text/plain")
		}

		fileData, err := ioutil.ReadFile("webfiles/" + strings.Replace(r.URL.String(), "/", "", 1))
		if err != nil {
			serverError(w, r, fmt.Sprintf("Could not read %v because of error: %v", r.URL, err))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(fileData)
	default:
		methodForbidden(w, r)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	indexFileData, err := ioutil.ReadFile("webfiles/html/index.html")
	if err != nil {
		serverError(w, r, fmt.Sprintf("Could not read 'webfiles/html/index.html' because of error: %v", err))
		return
	}

	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write(indexFileData)
	default:
		methodForbidden(w, r)
	}
}

func registerWebInterfaceHandlers() {
	http.HandleFunc("/", index)
	http.HandleFunc("/resources/", resources)
	http.HandleFunc("/favicon.ico", favicon)
}
