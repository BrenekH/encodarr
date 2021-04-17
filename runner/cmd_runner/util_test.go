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
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out, outErr := parseColonTimeToDuration(test.input)

			if test.errExpected && outErr == nil {
				t.Errorf("expected an error from %v", test.input)
			} else if test.errExpected && outErr != nil {
				return
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
