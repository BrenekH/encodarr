package cmd_runner

import "time"

// Sincer is an interface that allows mocking out time.Since for testing.
type Sincer interface {
	Since(time.Time) time.Duration
}
