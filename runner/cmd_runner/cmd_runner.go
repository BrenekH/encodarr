package cmd_runner

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

type CmdRunner struct {
	Executable string
	BaseArgs   []string
	done       bool
	fps        float64
	speed      float64
}

func (r *CmdRunner) Done() bool {
	return r.done
}

func (r *CmdRunner) Start(args []string) {
	a := append(r.BaseArgs, args...)
	c := exec.Command(r.Executable, a...)

	errPipe, _ := c.StderrPipe()
	b := make([]byte, 1024)

	fpsRe, err := regexp.Compile(`fps= *([0-9\.]*)`)
	if err != nil {
		panic(err)
	}

	timeRe, err := regexp.Compile(`time= *([0-9:\.]*)`)
	if err != nil {
		panic(err)
	}

	speedRe, err := regexp.Compile(`speed= *([0-9\.]*)`)
	if err != nil {
		panic(err)
	}

	go func() {
		fmt.Println("Starting FFmpeg command")
		c.Start()

		for {
			n, err := errPipe.Read(b)
			line := string(b[:n])
			fmt.Println(err, n)
			fmt.Println(line)

			// TODO: Move parsing to function
			// FPS
			r.fps, err = strconv.ParseFloat(fpsRe.FindStringSubmatch(line)[1], 64)
			if err != nil {
				panic(err)
			}

			// Time
			timeStr := timeRe.FindStringSubmatch(line)[1]

			// Speed
			r.speed, err = strconv.ParseFloat(speedRe.FindStringSubmatch(line)[1], 64)
			if err != nil {
				panic(err)
			}

			fmt.Println(r.fps, timeStr, r.speed)

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
