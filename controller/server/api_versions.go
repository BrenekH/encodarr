package server

import (
	"encoding/json"
	"net/http"

	"github.com/BrenekH/encodarr/controller/config"
)

var webVersions []string = make([]string, 1)
var runnerVersions []string = make([]string, 1)

type bothJSON struct {
	Web    singleJSON `json:"web"`
	Runner singleJSON `json:"runner"`
}

type singleJSON struct {
	Versions []string `json:"versions"`
}

func apiVersions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		resp := bothJSON{singleJSON{webVersions}, singleJSON{runnerVersions}}

		b, err := json.Marshal(resp)
		if err != nil {
			serverError(w, r, err.Error())
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	default:
		methodForbidden(w, r)
	}
}

func webAPIVersions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		b, err := json.Marshal(singleJSON{webVersions})
		if err != nil {
			serverError(w, r, err.Error())
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	default:
		methodForbidden(w, r)
	}
}

func runnerAPIVersions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		b, err := json.Marshal(singleJSON{runnerVersions})
		if err != nil {
			serverError(w, r, err.Error())
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	default:
		methodForbidden(w, r)
	}
}

func controllerVersion(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(config.Version))
}

func registerAPIHandlers() {
	registerWebAPIv1Handlers()
	webVersions[0] = "v1"

	registerRunnerAPIv1Handlers()
	runnerVersions[0] = "v1"

	http.HandleFunc("/api", apiVersions)
	http.HandleFunc("/api/web", webAPIVersions)
	http.HandleFunc("/api/runner", runnerAPIVersions)

	http.HandleFunc("/version", controllerVersion)
}
