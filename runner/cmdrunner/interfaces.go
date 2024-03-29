package cmdrunner

import (
	"io"
	"time"
)

// Sincer is an interface that allows mocking out time.Since for testing.
type Sincer interface {
	Since(t time.Time) time.Duration
}

// Commander is an interface that allows for mocking out the os/exec package for testing.
type Commander interface {
	Command(name string, args ...string) Cmder
}

// Cmder is an interface for mocking out the exec.Cmd struct.
type Cmder interface {
	Start() error
	StderrPipe() (io.ReadCloser, error)
	Wait() error
}
