package http

import (
	"io"
	netHTTP "net/http"
	"time"
)

// RequestDoer is a mock interface for the ApiV1 which allows
// for mock http Client's to be inserted during testing.
type RequestDoer interface {
	Do(*netHTTP.Request) (*netHTTP.Response, error)
}

// CurrentTimer is an interface for abstracting time.Now away from the code
// for testing purposes.
type CurrentTimer interface {
	Now() time.Time
}

// FSer is a mock interface for the ApiV1 struct which allows
// file operations to be mocked out during testing.
type FSer interface {
	Create(name string) (Filer, error)
	Open(name string) (Filer, error)
}

type Filer interface {
	io.Closer
	io.Reader
	io.Writer
	Name() string
}
