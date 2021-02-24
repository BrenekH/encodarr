package server

import "net/http"

// subRouter registers HTTP handlers with a BaseRoute prepended to the passed route
type subRouter struct {
	BaseRoute string
}

// HandleFunc registers the passed handler with the net/http package
func (r *subRouter) HandleFunc(subPattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(r.BaseRoute+subPattern, handler)
}

func newSubRouter(baseRoute string) subRouter {
	return subRouter{baseRoute}
}
