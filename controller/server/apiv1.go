package server

import (
	"net/http"
)

// subRouter registers HTTP handlers with a BaseRoute prepended to the passed route
type subRouter struct {
	BaseRoute string
}

func (r *subRouter) handleFunc(subPattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(r.BaseRoute+subPattern, handler)
}

func newSubRouter(baseRoute string) subRouter {
	return subRouter{baseRoute}
}

func apiSample(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><head><title>Project RedCedar - API Test</title></head><body><h4>Hello, World!</h4></body></html>`))
	default:
		methodForbidden(w, r)
	}
}

func registerAPIv1Handlers() {
	r := newSubRouter("/api/v1")

	r.handleFunc("/sample", apiSample)
}
