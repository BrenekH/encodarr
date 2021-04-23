package runner

import "errors"

// ErrUnresponsive represents the error state when the Controller decides that the Runner is no longer responsive.
var ErrUnresponsive error = errors.New("received unresponsive status code")
