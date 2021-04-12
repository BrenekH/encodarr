package cmd_runner

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

var (
	fpsRe   *regexp.Regexp
	timeRe  *regexp.Regexp
	speedRe *regexp.Regexp
)

func init() {
	// Create Regexes and panic if they fail to compile.
	// This allows them to be used for every regex search instead of needing to recompile everytime.
	var err error
	fpsRe, err = regexp.Compile(`fps= *([0-9\.]*)`)
	if err != nil {
		panic(err)
	}

	timeRe, err = regexp.Compile(`time= *([0-9:\.]*)`)
	if err != nil {
		panic(err)
	}

	speedRe, err = regexp.Compile(`speed= *([0-9\.]*)`)
	if err != nil {
		panic(err)
	}
}

type CmdRunner struct {
	Executable string
	BaseArgs   []string
	done       bool
	fps        float64
	time       string
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

	go func() {
		fmt.Println("Starting FFmpeg command")
		c.Start()

		for {
			n, err := errPipe.Read(b)
			line := string(b[:n])
			fmt.Println(err, n)
			fmt.Println(line)

			fps, time, speed := parseFFmpegLine(line)

			r.fps = fps
			r.time = time
			r.speed = speed

			fmt.Println(r.fps, r.time, r.speed)

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

// parseFFmpegLine parses out the fps, time, and speed information from a standard FFmpeg statistics line.
// There might be a speed up that involves changing the line parameter(and maybe the return results) to a pointer
// (avoids copying the value for a new frame), but the jury is still out on that one.
func parseFFmpegLine(line string) (fps float64, time string, speed float64) {
	// FPS
	fps, err := strconv.ParseFloat(fpsRe.FindStringSubmatch(line)[1], 64)
	if err != nil {
		panic(err)
	}

	// Time
	time = timeRe.FindStringSubmatch(line)[1]

	// Speed
	speed, err = strconv.ParseFloat(speedRe.FindStringSubmatch(line)[1], 64)
	if err != nil {
		panic(err)
	}
	return
}
