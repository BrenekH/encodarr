package cmd_runner

import "time"

type TimeSince struct{}

func (s TimeSince) Since(t time.Time) time.Duration {
	return time.Since(t)
}
