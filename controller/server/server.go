package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func serverError(w http.ResponseWriter, r *http.Request, reason string) {
	fmt.Println(reason)
	w.WriteHeader(500)
	w.Write([]byte(`<html><body>Code 500: Server Error</body></html>`))
}

func favicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/x-icon")

	icoData, err := ioutil.ReadFile("webfiles/favicon.ico")
	if err != nil {
		serverError(w, r, fmt.Sprintf("Could not read %v because of error: %v", r.URL, err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(icoData)
}

// resources is a http handler for resource(css, js) requests to the server.
func resources(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Addr: %v; Page: %v; Time: %v; X-Forwarded_For: %v\n", r.RemoteAddr, r.URL, time.Now().String(), r.Header.Get("x-forwarded-for"))
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
			fmt.Printf("Could not identify MIME type for resources request: %v\n", r.URL.String())
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
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`<html><body><h1>Project RedCedar</h1><p>METHOD FORBIDDEN: That HTTP method is not allowed to this route</p></body></html>`))
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	indexFileData, err := ioutil.ReadFile("webfiles/html/index.html")
	if err != nil {
		serverError(w, r, fmt.Sprintf("Could not read 'webfiles/html/index.html' because of error: %v", err))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write(indexFileData)
	default:
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`<html><body><h1>Project RedCedar</h1><p>METHOD FORBIDDEN: That HTTP method is not allowed to this route</p></body></html>`))
	}
}

// RunHTTPServer runs the HTTP server for Project RedCedar in a blocking manner.
func RunHTTPServer() {
	envPort := "8080"
	http.Handle("/", http.HandlerFunc(index))
	http.HandleFunc("/resources/", resources)
	http.HandleFunc("/favicon.ico", favicon)
	log.Fatal(http.ListenAndServe(":"+envPort, nil))
}
