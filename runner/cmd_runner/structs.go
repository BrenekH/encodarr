package cmd_runner

import (
	"os/exec"
	"time"
)

type TimeSince struct{}

func (s TimeSince) Since(t time.Time) time.Duration {
	return time.Since(t)
}

type ExecCommander struct{}

func (e ExecCommander) Command(name string, args ...string) Cmder {
	return exec.Command(name, args...)
}
