package options

import (
	"os"
	"testing"
)

func TestStringVarFromEnv(t *testing.T) {
	t.Run("Load Valid EnvVar", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("MY_VAR", "hello")
		s := "default"
		stringVarFromEnv(&s, "MY_VAR")

		if s != "hello" {
			t.Errorf("expected 'hello' but got '%v'", s)
		}
	})

	t.Run("Load Invalid EnvVar", func(t *testing.T) {
		os.Clearenv()
		s := "default"
		stringVarFromEnv(&s, "MY_VAR")

		if s != "default" {
			t.Errorf("expected 'default' but got '%v'", s)
		}
	})
}

func TestConfigDir(t *testing.T) {
	c := ConfigDir()
	v := configDir

	if c != v {
		t.Errorf("expected %v but got %v", v, c)
	}
}

func TestTempDir(t *testing.T) {
	td := TempDir()
	v := tempDir

	if td != v {
		t.Errorf("expected %v but got %v", v, td)
	}
}

func TestRunnerName(t *testing.T) {
	r := RunnerName()
	v := runnerName

	if r != v {
		t.Errorf("expected %v but got %v", v, r)
	}
}

func TestControllerIP(t *testing.T) {
	cip := ControllerIP()
	v := controllerIP

	if cip != v {
		t.Errorf("expected %v but got %v", v, cip)
	}
}

func TestControllerPort(t *testing.T) {
	port := ControllerPort()
	v := controllerPort

	if port != v {
		t.Errorf("expected %v but got %v", v, port)
	}
}

func TestInTestMode(t *testing.T) {
	tm := InTestMode()
	v := inTestMode

	if tm != v {
		t.Errorf("expected %v but got %v", v, tm)
	}
}
