package cmdrunner

import (
	"os/exec"
	"time"
)

type timeSince struct{}

func (s timeSince) Since(t time.Time) time.Duration {
	return time.Since(t)
}

type execCommander struct{}

func (e execCommander) Command(name string, args ...string) Cmder {
	return exec.Command(name, args...)
}
