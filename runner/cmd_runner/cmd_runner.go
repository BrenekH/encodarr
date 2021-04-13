package cmd_runner

// TODO: Remove panics in favor of actual logging.

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"github.com/BrenekH/encodarr/runner"
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
	Executable   string
	BaseArgs     []string
	fileDuration time.Duration
	startTime    time.Time
	done         bool
	fps          float64
	time         string
	speed        float64
}

func (r *CmdRunner) Done() bool {
	return r.done
}

func (r *CmdRunner) Start(ji runner.JobInfo) {
	dur, err := strconv.ParseInt(ji.MediaInfo.General.Duration, 10, 64)
	if err != nil {
		panic(err)
	}
	// ji.MediaInfo.General.Duration is in milliseconds
	r.fileDuration = time.Duration(dur) * time.Millisecond

	a := append(r.BaseArgs, ji.CommandArgs...)
	c := exec.Command(r.Executable, a...)

	errPipe, _ := c.StderrPipe()
	b := make([]byte, 1024)

	r.startTime = time.Now()

	go func() {
		fmt.Println("Starting FFmpeg command")
		c.Start()

		for {
			n, _ := errPipe.Read(b)
			line := string(b[:n])
			// fmt.Println(err, n)
			// fmt.Println(line)

			fps, time, speed := parseFFmpegLine(line)

			r.fps = fps
			r.time = time
			r.speed = speed

			// fmt.Println(r.fps, r.time, r.speed)

			if n == 0 {
				break
			}
		}

		c.Wait()

		r.done = true
		fmt.Println("FFmpeg command finished")
	}()
}

func (r *CmdRunner) Status() runner.JobStatus {
	currentFileTime, err := parseFFmpegTime(r.time)
	if err != nil {
		currentFileTime = time.Duration(0)
	}

	return runner.JobStatus{
		Stage:                       "Running FFmpeg",
		Percentage:                  fmt.Sprintf("%.2f", (float64(currentFileTime)/float64(r.fileDuration))*100),
		JobElapsedTime:              fmt.Sprintf("%v", time.Since(r.startTime).Round(time.Second).String()),
		FPS:                         fmt.Sprintf("%v", r.fps),
		StageElapsedTime:            fmt.Sprintf("%v", time.Since(r.startTime).Round(time.Second).String()),
		StageEstimatedTimeRemaining: fmt.Sprintf("%v", time.Duration(float64(r.fileDuration-currentFileTime)/r.speed).Round(time.Second).String()),
	}
}

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
	var err error

	// FPS
	fpsReMatch := fpsRe.FindStringSubmatch(line)
	if len(fpsReMatch) > 1 {
		fps, err = strconv.ParseFloat(fpsReMatch[1], 64)
		if err != nil {
			panic(err)
		}
	}

	// Time
	timeReMatch := timeRe.FindStringSubmatch(line)
	if len(timeReMatch) > 1 {
		time = timeReMatch[1]
	}

	// Speed
	speedReMatch := speedRe.FindStringSubmatch(line)
	if len(speedReMatch) > 1 {
		speed, err = strconv.ParseFloat(speedReMatch[1], 64)
		if err != nil {
			fmt.Println(err)
		}
	}
	return
}

// parseFFmpegTime takes a "HH:MM:SS" and converts it to a time.Duration.
// The hour portion does not have to be <= 24.
func parseFFmpegTime(s string) (time.Duration, error) {
	var hrs, mins, secs int64

	_, err := fmt.Sscanf(s, "%d:%d:%d", &hrs, &mins, &secs)
	if err != nil {
		return time.Duration(0), err
	}

	// Making everything into time.Durations probably isn't the best option,
	// but there doesn't seem to be a great option to avoid them and still return a time.Duration.
	return time.Duration(
		time.Hour*time.Duration(hrs) +
			time.Minute*time.Duration(mins) +
			time.Second*time.Duration(secs),
	), nil
}
