package cmd_runner

import (
	"reflect"
	"testing"
	"time"

	"github.com/BrenekH/encodarr/runner"
)

func TestCmdRunnerDone(t *testing.T) {
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

func TestCmdRunnerStart(t *testing.T) {
	// This test won't be implemented until I figure out
	// how to prevent a command from actually being ran.
}

func TestCmdRunnerStatus(t *testing.T) {
	tests := []struct {
		name         string
		fps          float64
		time         string
		speed        float64
		startTime    time.Time
		fileDuration time.Duration
		expected     runner.JobStatus
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := CmdRunner{
				fps:          test.fps,
				time:         test.time,
				speed:        test.speed,
				startTime:    test.startTime,
				fileDuration: test.fileDuration,
			}

			out := r.Status()

			if !reflect.DeepEqual(out, test.expected) {
				t.Errorf("expected %v, but got %v", test.expected, out)
			}
		})
	}
}

func TestCmdRunnerResults(t *testing.T) {
	tests := []struct {
		name      string
		failed    bool
		errors    []string
		warnings  []string
		startTime time.Time
		expected  runner.CommandResults
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Set all the private vars, relating to the creation of the runner.CommandResults struct.
			// Check if the generated struct is equal to the expected one.
			r := CmdRunner{
				failed:    test.failed,
				errors:    test.errors,
				warnings:  test.warnings,
				startTime: test.startTime,
			}

			out := r.Results()

			if !reflect.DeepEqual(out, test.expected) {
				t.Errorf("expected %v, but got %v", test.expected, out)
			}
		})
	}
}
