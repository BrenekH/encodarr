package cmd_runner

import (
	"reflect"
	"testing"
	"time"
)

func TestParseColonTimeToDuration(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    time.Duration
		errExpected bool
	}{
		{name: "Zeroes", input: "00:00:00", expected: time.Duration(0), errExpected: false},

		{name: "Seconds", input: "00:00:30", expected: time.Duration(30 * time.Second), errExpected: false},
		{name: "Minutes", input: "00:45:00", expected: time.Duration(45 * time.Minute), errExpected: false},
		{name: "Hours", input: "10:00:00", expected: time.Duration(10 * time.Hour), errExpected: false},

		{name: "Over 60 Seconds", input: "00:00:99", expected: time.Duration(99 * time.Second), errExpected: false},
		{name: "Over 60 Minutes", input: "00:71:00", expected: time.Duration(71 * time.Minute), errExpected: false},
		{name: "Over 24 Hours", input: "60:00:00", expected: time.Duration(60 * time.Hour), errExpected: false},

		{name: "Combined", input: "03:17:56", expected: time.Duration(3*time.Hour + 17*time.Minute + 56*time.Second), errExpected: false},
		{name: "3 Digit Times", input: "100:200:300", expected: time.Duration(100*time.Hour + 200*time.Minute + 300*time.Second), errExpected: false},

		{name: "Missing \"Hour\" Spot", input: "35:00", expected: time.Duration(35 * time.Minute), errExpected: true},
		{name: "Extra Chars", input: "hi10:10:10", expected: time.Duration(10*time.Hour + 10*time.Minute + 10*time.Second), errExpected: true},
		{name: "Resembles Valid Input", input: "HH:MM:SS", errExpected: true},
		{name: "Complete Nonsense", input: "Hello, World!", errExpected: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out, outErr := parseColonTimeToDuration(test.input)

			if test.errExpected && outErr == nil {
				t.Errorf("expected an error from %v", test.input)
			} else if test.errExpected && outErr != nil {
				return
			} else if !test.errExpected && outErr != nil {
				t.Errorf("unexpected error '%v' from %v", outErr, test.input)
			}

			if !reflect.DeepEqual(out, test.expected) {
				t.Errorf("got %v, expected %v", out, test.expected)
			}
		})
	}
}

func TestParseFFmpegLine(t *testing.T) {
	tests := []struct {
		name     string
		inLine   string
		outFps   float64
		outTime  string
		outSpeed float64
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				tFps   float64
				tTime  string
				tSpeed float64
			)

			parseFFmpegLine(test.inLine, &tFps, &tTime, &tSpeed)

			if tFps != test.outFps {
				t.Errorf("fps: got %v, expected %v", tFps, test.outFps)
			}
			if tTime != test.outTime {
				t.Errorf("time: got %v, expected %v", tTime, test.outTime)
			}
			if tSpeed != test.outSpeed {
				t.Errorf("speed: got %v, expected %v", tSpeed, test.outSpeed)
			}
		})
	}
}

func TestExtractFps(t *testing.T) {
	tests := []struct {
		name        string
		inLine      string
		expected    float64
		errExpected bool
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out, outErr := extractFps(test.inLine)

			if test.errExpected && outErr == nil {
				t.Errorf("expected an error from %v", test.inLine)
			} else if test.errExpected && outErr != nil {
				return
			}

			if out != test.expected {
				t.Errorf("got %v, expected %v", out, test.expected)
			}
		})
	}
}

func TestExtractTime(t *testing.T) {
	tests := []struct {
		name        string
		inLine      string
		expected    string
		errExpected bool
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out, outErr := extractTime(test.inLine)

			if test.errExpected && outErr == nil {
				t.Errorf("expected an error from %v", test.inLine)
			} else if test.errExpected && outErr != nil {
				return
			}

			if out != test.expected {
				t.Errorf("got %v, expected %v", out, test.expected)
			}
		})
	}
}

func TestExtractSpeed(t *testing.T) {
	tests := []struct {
		name        string
		inLine      string
		expected    float64
		errExpected bool
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out, outErr := extractSpeed(test.inLine)

			if test.errExpected && outErr == nil {
				t.Errorf("expected an error from %v", test.inLine)
			} else if test.errExpected && outErr != nil {
				return
			}

			if out != test.expected {
				t.Errorf("got %v, expected %v", out, test.expected)
			}
		})
	}
}
