package options

import (
	"reflect"
	"testing"
)

func TestStringVar(t *testing.T) {
	t.Run("Add Flag to Flags Slice", func(t *testing.T) {
		flags = []flagger{}

		var s string

		stringVar(&s, "test", "")
		expected := []flagger{StringFlag{name: "test", usage: "", pointer: &s}}

		if !reflect.DeepEqual(flags, expected) {
			t.Errorf("expected %v but got %v", expected, flags)
		}
	})
}

func TestStringFlag(t *testing.T) {
	s := "default"
	sF := StringFlag{
		name:    "test",
		usage:   "Use test",
		pointer: &s,
	}

	t.Run("Get Name", func(t *testing.T) {
		if sF.Name() != sF.name {
			t.Errorf("expected '%v' but got '%v", sF.name, sF.Name())
		}
	})

	t.Run("Get Usage", func(t *testing.T) {
		if sF.Usage() != sF.usage {
			t.Errorf("expected '%v' but got '%v", sF.usage, sF.Usage())
		}
	})

	t.Run("Set string", func(t *testing.T) {
		outErr := sF.Parse("hello")

		if outErr != nil {
			t.Errorf("unexpected error: %v", outErr)
		}

		if s != "hello" {
			t.Errorf("expected 'hello' but got '%v", s)
		}
	})
}
