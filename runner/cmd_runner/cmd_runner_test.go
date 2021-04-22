package cmd_runner

import (
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
	// Case: done, failed, warnings, errors are reset when the command starts
	// Case: Commander.Command is called with appropriate args
}

func TestStartResults(t *testing.T) {
	// Case: Zero exit code causes Results to set Failed to false
	// Case: Non-zero exit code causes Results to set Failed to true
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
		expected  runner.CommandResults
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cR := NewCmdRunner()

			cR.failed = test.failed
			cR.startTime = test.startTime
			cR.warnings = test.warnings
			cR.errors = test.errors

			results := cR.Results()

			if !reflect.DeepEqual(results, test.expected) {
				t.Errorf("expected %v but got %v", test.expected, results)
			}
		})
	}
}
