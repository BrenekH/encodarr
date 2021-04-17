package cmd_runner

import (
	"testing"
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
