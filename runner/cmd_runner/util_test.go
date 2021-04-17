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

		// Errors
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
	}{
		{name: "Minimal Data", inLine: "fps= 4.5 time= 08:12:03 speed= 0.420", outFps: 4.5, outTime: "08:12:03", outSpeed: 0.42},
		{name: "Full FFmpeg Line", inLine: "frame=  105 fps= 28 q=28.0 size=     256kB time=00:00:04.30 bitrate= 486.7kbits/s dup=22 drop=0 speed=1.17x",
			outFps:   28.0,
			outTime:  "00:00:04.30",
			outSpeed: 1.17,
		},

		{name: "Only Fps", inLine: "fps= 4.5 ", outFps: 4.5, outTime: "", outSpeed: 0.0},
		{name: "Only Time", inLine: "time= 08:12:03 ", outFps: 0.0, outTime: "08:12:03", outSpeed: 0.0},
		{name: "Only Speed", inLine: "speed= 0.420", outFps: 0.0, outTime: "", outSpeed: .42},

		{name: "No Actionable Data", inLine: "Hello, World!", outFps: 0.0, outTime: "", outSpeed: 0.0},
	}

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

	// This test can't be run with the general loop because we can't
	// set the value of a var passed into parseFFmpegLine in that system.
	// I could add that ability, but I'd rather not tailor the entire struct
	// just for one test.
	t.Run("Zero Speed Doesn't Modify Value", func(t *testing.T) {
		var (
			tFps   float64
			tTime  string
			tSpeed float64
		)

		tSpeed = 1.0

		parseFFmpegLine("speed= 0.0", &tFps, &tTime, &tSpeed)

		if tSpeed != 1.0 {
			t.Errorf("expected speed to be 1.0 but instead it was %v", tSpeed)
		}
	})
}

func TestExtractFps(t *testing.T) {
	tests := []struct {
		name        string
		inLine      string
		expected    float64
		errExpected bool
	}{
		{name: "Basic", inLine: "fps= 500.1 ", expected: 500.1, errExpected: false},
		{name: "No Spaces", inLine: "fps=10.2 ", expected: 10.2, errExpected: false},
		{name: "Absurd Amount of Spaces", inLine: "fps=                10.2 ", expected: 10.2, errExpected: false},
		{name: "Full Line", inLine: "frame=  105 fps= 28 q=28.0 size=     256kB time=00:00:04.30 bitrate= 486.7kbits/s dup=22 drop=0 speed=1.17x",
			expected: 28.0, errExpected: false},

		// Errors
		{name: "Missing Space", inLine: "fps= 23.09", expected: 0.0, errExpected: true},
		{name: "Is a String", inLine: "fps= hello", expected: 0.0, errExpected: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out, outErr := extractFps(test.inLine)

			if test.errExpected && outErr == nil {
				t.Errorf("expected an error from %v", test.inLine)
			} else if test.errExpected && outErr != nil {
				return
			} else if !test.errExpected && outErr != nil {
				t.Errorf("unexpected error '%v' from %v", outErr, test.inLine)
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
			} else if !test.errExpected && outErr != nil {
				t.Errorf("unexpected error '%v' from %v", outErr, test.inLine)
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
			} else if !test.errExpected && outErr != nil {
				t.Errorf("unexpected error '%v' from %v", outErr, test.inLine)
			}

			if out != test.expected {
				t.Errorf("got %v, expected %v", out, test.expected)
			}
		})
	}
}
