package cmd_runner

import (
	"fmt"
	"os/exec"
)

type CmdRunner struct {
	Executable string
	BaseArgs   []string
	done       bool
}

func (r *CmdRunner) Done() bool {
	return r.done
}

func (r *CmdRunner) Start(args []string) {
	a := append(r.BaseArgs, args...)
	c := exec.Command(r.Executable, a...)

	errPipe, _ := c.StderrPipe()
	b := make([]byte, 1024)

	go func() {
		fmt.Println("Starting FFmpeg command")
		c.Start()

		for {
			n, err := errPipe.Read(b)
			fmt.Println(err, n)
			fmt.Println(string(b[:n]))
			// TODO: Parse out status from line

			if n == 0 {
				break
			}
		}

		r.done = true
		fmt.Println("FFmpeg command finished")
	}()
}

func (r *CmdRunner) Status() {}

func NewCmdRunner() CmdRunner {
	return CmdRunner{
		Executable: "ffmpeg",
		BaseArgs:   []string{"-hide_banner", "-loglevel", "warning", "-stats", "-y"},
		done:       false,
	}
}
