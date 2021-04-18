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
