package cmd_runner

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"github.com/BrenekH/encodarr/runner"
	"github.com/BrenekH/logange"
)

var (
	fpsRe   *regexp.Regexp
	timeRe  *regexp.Regexp
	speedRe *regexp.Regexp

	logger logange.Logger
)

func init() {
	logger = logange.NewLogger("cmd_runner")

	// Create Regexes and exit (logger.Critical) if they fail to compile.
	// This allows them to be used for every regex search instead of needing to recompile everytime.
	var err error
	fpsRe, err = regexp.Compile(`fps= *([0-9\.]*)`)
	if err != nil {
		logger.Critical(err.Error())
	}

	timeRe, err = regexp.Compile(`time= *([0-9:\.]*)`)
	if err != nil {
		logger.Critical(err.Error())
	}

	speedRe, err = regexp.Compile(`speed= *([0-9\.]*)`)
	if err != nil {
		logger.Critical(err.Error())
	}
}

type CmdRunner struct {
	Executable   string
	BaseArgs     []string
	fileDuration time.Duration
	startTime    time.Time
	done         bool
	failed       bool
	warnings     []string
	errors       []string
	fps          float64
	time         string
	speed        float64
}

func (r *CmdRunner) Done() bool {
	return r.done
}

func (r *CmdRunner) Start(ji runner.JobInfo) {
	// These variables need to be reset on every run because they only apply to one run,
	// but the CmdRunner persists over many command runs.
	r.done = false
	r.failed = false
	r.warnings = []string{}
	r.errors = []string{}

	dur, err := strconv.ParseInt(ji.MediaInfo.General.Duration, 10, 64)
	if err != nil {
		logger.Error(err.Error())
		r.errors = append(r.errors, err.Error())
		r.done = true
		r.failed = true
		return
	}

	// ji.MediaInfo.General.Duration is in milliseconds
	r.fileDuration = time.Duration(dur) * time.Millisecond

	a := append(r.BaseArgs, ji.CommandArgs...)
	c := exec.Command(r.Executable, a...)

	errPipe, _ := c.StderrPipe()
	b := make([]byte, 1024)

	r.startTime = time.Now()

	go func() {
		logger.Info("Starting FFmpeg command")
		c.Start()

		for {
			n, _ := errPipe.Read(b)
			line := string(b[:n])

			logger.Trace(fmt.Sprintf("%v %v", err, n))
			logger.Trace(line)

			parseFFmpegLine(line, &r.fps, &r.time, &r.speed)

			if n == 0 {
				break
			}
		}

		err = c.Wait()
		if err != nil {
			r.failed = true
			if exiterr, ok := err.(*exec.ExitError); ok {
				r.errors = append(r.errors, fmt.Sprintf("FFmpeg returned exit code: %v", exiterr.ExitCode()))
			} else {
				r.errors = append(r.errors, err.Error())
			}
		}

		r.done = true
		logger.Info("FFmpeg command finished")
	}()
}

func (r *CmdRunner) Status() runner.JobStatus {
	currentFileTime, err := parseColonTimeToDuration(r.time)
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

func (r *CmdRunner) Results() runner.CommandResults {
	return runner.CommandResults{
		Failed:         r.failed,
		JobElapsedTime: time.Since(r.startTime).Round(time.Second),
		Warnings:       r.warnings,
		Errors:         r.errors,
	}
}

func NewCmdRunner() CmdRunner {
	return CmdRunner{
		Executable: "ffmpeg",
		BaseArgs:   []string{"-hide_banner", "-loglevel", "warning", "-stats", "-y"},
	}
}

// parseFFmpegLine parses the fps, time, and speed information from a standard FFmpeg statistics line
// and updates the provided pointers if the parsing doesn't return an error.
//lint:ignore SA4009 The purpose of this code is to set and not read.
func parseFFmpegLine(line string, fps *float64, time *string, speed *float64) {
	// FPS
	if pFps, err := extractFps(line); err != nil {
		fps = &pFps
	} else {
		logger.Trace(err.Error())
	}

	// Time
	if pTime, err := extractTime(line); err != nil {
		time = &pTime
	} else {
		logger.Trace(err.Error())
	}

	// Speed
	if pSpeed, err := extractSpeed(line); err != nil {
		speed = &pSpeed
	} else {
		logger.Trace(err.Error())
	}
}

func extractFps(line string) (fps float64, err error) {
	fpsReMatch := fpsRe.FindStringSubmatch(line)
	if len(fpsReMatch) > 1 {
		fps, err = strconv.ParseFloat(fpsReMatch[1], 64)
		if err != nil {
			return
		}
	} else {
		err = fmt.Errorf("extractFps: fps regex returned too little groups (need at least 2)")
		return
	}
	return
}

func extractTime(line string) (time string, err error) {
	timeReMatch := timeRe.FindStringSubmatch(line)
	if len(timeReMatch) > 1 {
		time = timeReMatch[1]
	} else {
		err = fmt.Errorf("extractTime: time regex returned too little groups (need at least 2)")
		return
	}
	return
}

func extractSpeed(line string) (speed float64, err error) {
	speedReMatch := speedRe.FindStringSubmatch(line)
	if len(speedReMatch) > 1 {
		speed, err = strconv.ParseFloat(speedReMatch[1], 64)
		if err != nil {
			return
		}
	} else {
		err = fmt.Errorf("extractSpeed: speed regex returned too little groups (need at least 2)")
		return
	}
	return
}

// parseColonTimeToDuration takes a "HH:MM:SS" and converts it to a time.Duration.
// The hour portion does not have to be <= 24.
func parseColonTimeToDuration(s string) (time.Duration, error) {
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
