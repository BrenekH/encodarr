package http

import netHTTP "net/http"

// RequestDoer is a mock interface for the ApiV1 which allows
// for mock http Client's to be inserted during testing.
type RequestDoer interface {
	Do(*netHTTP.Request) (*netHTTP.Response, error)
}
