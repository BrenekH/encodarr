package cmd_runner

import (
	"io"
	"os/exec"
	"time"
)

type TimeSince struct{}

func (s TimeSince) Since(t time.Time) time.Duration {
	return time.Since(t)
}

type ExecCommander struct{}

func (e ExecCommander) Command(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}

type ExecCmder struct {
	internalCmd *exec.Cmd
}

func (e *ExecCmder) Start() error {
	return e.internalCmd.Start()
}

func (e *ExecCmder) StderrPipe() (io.ReadCloser, error) {
	return e.internalCmd.StderrPipe()
}

func (e *ExecCmder) Wait() error {
	return e.internalCmd.Wait()
}
