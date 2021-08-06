package cmdrunner

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/BrenekH/encodarr/runner"
)

func TestDone(t *testing.T) {
	cR := NewCmdRunner()

	t.Run("Done is true", func(t *testing.T) {
		cR.done = true

		if !cR.Done() {
			t.Errorf("expected true but got false")
		}
	})

	t.Run("Done is false", func(t *testing.T) {
		cR.done = false

		if cR.Done() {
			t.Errorf("expected false but got true")
		}
	})
}

func TestStart(t *testing.T) {
	t.Run("Critical Per-Job Variables are Reset", func(t *testing.T) {
		cR := NewCmdRunner()
		cR.cmdr = &mockCommander{}

		cR.done = true
		cR.failed = true
		cR.warnings = []string{"hello", "test"}
		cR.errors = []string{"hello", "test"}

		cR.Start(runner.JobInfo{})

		if cR.done {
			t.Errorf("done was supposed to be reset to false")
		}

		if cR.failed {
			t.Errorf("failed was supposed to be reset to false")
		}

		if !reflect.DeepEqual(cR.warnings, []string{}) {
			t.Errorf("warnings was supposed to be reset to an empty slice")
		}

		if !reflect.DeepEqual(cR.errors, []string{}) {
			t.Errorf("errors was supposed to be reset to an empty slice")
		}
	})

	t.Run("Check Args Passed to Commander.Command", func(t *testing.T) {
		mCmdr := mockCommander{}
		cR := NewCmdRunner()
		cR.cmdr = &mCmdr

		cR.Start(runner.JobInfo{
			CommandArgs: []string{"-i", "input.mp4", "output.mkv"},
		})

		expected := []string{"-hide_banner", "-loglevel", "warning", "-stats", "-y", "-i", "input.mp4", "output.mkv"}
		if !reflect.DeepEqual(mCmdr.lastCallArgs, expected) {
			t.Errorf("expected Commander.Command to be called with %v, but got %v instead", expected, mCmdr)
		}
	})
}

func TestStartResults(t *testing.T) {
	t.Run("Exit Code 0, Failed is false", func(t *testing.T) {
		cR := NewCmdRunner()
		cR.cmdr = &mockCommander{cmder: mockCmder{statusCode: 0}}

		cR.Start(runner.JobInfo{})

		results := cR.Results()
		if results.Failed {
			t.Errorf("expected failed to be false, not true")
		}
	})

	t.Run("Exit Code 1, Failed is true", func(t *testing.T) {
		cR := NewCmdRunner()
		cR.cmdr = &mockCommander{cmder: mockCmder{statusCode: 1}}

		cR.Start(runner.JobInfo{})

		// Because there is a goroutine inside cR.Start that has to finish before the results are ready,
		// we just get to wait for it to be complete.
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
		defer cancel()
		if failedToSetDone := waitCmdRunner(&ctx, &cR); failedToSetDone {
			t.Errorf("test CmdRunner failed to set the done variable within 10 seconds")
		}

		results := cR.Results()
		if !results.Failed {
			t.Errorf("expected failed to be true, not false")
		}
	})
}

func TestStatus(t *testing.T) {
	tests := []struct {
		name         string
		fileDuration time.Duration
		startTime    time.Time
		fps          float64
		time         string
		speed        float64

		timeNow  time.Time
		expected runner.JobStatus
	}{
		{
			name:    "Nothing Custom Set",
			timeNow: time.Unix(0, 0).UTC(),
			expected: runner.JobStatus{
				Stage:                       "Running FFmpeg",
				Percentage:                  "NaN",
				JobElapsedTime:              "2562047h47m16.854775807s",
				FPS:                         "0",
				StageElapsedTime:            "2562047h47m16.854775807s",
				StageEstimatedTimeRemaining: "-2562047h47m16.854775808s",
			},
		},
		{
			name:         "Everything Custom Set",
			timeNow:      time.Date(2000, time.January, 1, 0, 20, 0, 0, time.UTC),
			fileDuration: time.Duration(20) * time.Minute,
			startTime:    time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
			fps:          9.001,
			time:         "00:10:00",
			speed:        1.0,
			expected: runner.JobStatus{
				Stage:                       "Running FFmpeg",
				Percentage:                  "50.00",
				JobElapsedTime:              "20m0s",
				FPS:                         "9.001",
				StageElapsedTime:            "20m0s",
				StageEstimatedTimeRemaining: "10m0s",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cR := NewCmdRunner()

			cR.fileDuration = test.fileDuration
			cR.startTime = test.startTime
			cR.fps = test.fps
			cR.time = test.time
			cR.speed = test.speed

			cR.timeSince = &mockSincer{t: test.timeNow}

			js := cR.Status()

			if !reflect.DeepEqual(js, test.expected) {
				t.Errorf("expected %v but got %v", test.expected, js)
			}
		})
	}
}

func TestResults(t *testing.T) {
	tests := []struct {
		name      string
		failed    bool
		startTime time.Time
		warnings  []string
		errors    []string

		timeNow  time.Time
		expected runner.CommandResults
	}{
		{
			name:      "No Warnings or Errors",
			failed:    false,
			startTime: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
			warnings:  []string{},
			errors:    []string{},
			timeNow:   time.Date(2000, time.January, 1, 0, 20, 0, 0, time.UTC),
			expected: runner.CommandResults{
				Failed:         false,
				JobElapsedTime: time.Duration(20) * time.Minute,
				Warnings:       []string{},
				Errors:         []string{},
			},
		},
		{
			name:      "Only Warnings (failed = false)",
			failed:    false,
			startTime: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
			warnings:  []string{"Unsupported subtitle codec for container: mkv"},
			errors:    []string{},
			timeNow:   time.Date(2000, time.January, 1, 0, 20, 0, 0, time.UTC),
			expected: runner.CommandResults{
				Failed:         false,
				JobElapsedTime: time.Duration(20) * time.Minute,
				Warnings:       []string{"Unsupported subtitle codec for container: mkv"},
				Errors:         []string{},
			},
		},
		{
			name:      "Job Failed (failed = true, errors)",
			failed:    true,
			startTime: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
			warnings:  []string{},
			errors:    []string{"FFmpeg returned a non-zero exit code: 1"},
			timeNow:   time.Date(2000, time.January, 1, 0, 20, 0, 0, time.UTC),
			expected: runner.CommandResults{
				Failed:         true,
				JobElapsedTime: time.Duration(20) * time.Minute,
				Warnings:       []string{},
				Errors:         []string{"FFmpeg returned a non-zero exit code: 1"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cR := NewCmdRunner()

			cR.failed = test.failed
			cR.startTime = test.startTime
			cR.warnings = test.warnings
			cR.errors = test.errors

			cR.timeSince = &mockSincer{t: test.timeNow}

			results := cR.Results()

			if !reflect.DeepEqual(results, test.expected) {
				t.Errorf("expected %v but got %v", test.expected, results)
			}
		})
	}
}

// waitCmdRunner is a helper function that simply waits until the passed
// CmdRunner indicates it is done or the context is up.
//
// The boolean return value indicates whether or not the CmdRunner failed to set the done variable.
func waitCmdRunner(ctx *context.Context, cR *CmdRunner) bool {
	for {
		if runner.IsContextFinished(ctx) {
			return true
		}
		if cR.Done() {
			return false
		}
	}
}
