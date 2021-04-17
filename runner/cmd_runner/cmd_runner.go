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

// The regexes are defined here instead of util.go,
// so that their instantiation in init() can be in the same file as the var declerations.
var (
	fpsRe   *regexp.Regexp
	timeRe  *regexp.Regexp
	speedRe *regexp.Regexp

	logger logange.Logger
)

// To avoid init function run order confusion I only want one in this package.
func init() {
	logger = logange.NewLogger("cmd_runner")

	// Create Regexes and exit (logger.Critical) if they fail to compile.
	// This allows them to be used for every regex search instead of needing to recompile everytime.
	var err error
	fpsRe, err = regexp.Compile(`fps= *([0-9\.]*) `)
	if err != nil {
		logger.Critical(err.Error())
	}

	// This regex is particularly specific so that false positives don't cause the UI to screw up.
	timeRe, err = regexp.Compile(`time= *([0-9\.]*:[0-9\.]*:[0-9\.]*) `)
	if err != nil {
		logger.Critical(err.Error())
	}

	speedRe, err = regexp.Compile(`speed= *([0-9\.]*)`)
	if err != nil {
		logger.Critical(err.Error())
	}
}

func NewCmdRunner() CmdRunner {
	return CmdRunner{
		Executable: "ffmpeg",
		BaseArgs:   []string{"-hide_banner", "-loglevel", "warning", "-stats", "-y"},
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
