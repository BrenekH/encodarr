package http

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseFFmpegCmd(t *testing.T) {
	tests := []struct {
		name     string
		inFname  string
		outFname string
		inCmd    []string
		expected []string
	}{
		{
			name:     "Encode to HEVC",
			inFname:  "input.mkv",
			outFname: "output.mkv",
			inCmd:    []string{"-i", "ENCODARR_INPUT_FILE", "-map", "0:s?", "-map", "0:a", "-c", "copy", "-map", "0:v", "-vcodec", "hevc"},
			expected: []string{"-i", "input.mkv", "-map", "0:s?", "-map", "0:a", "-c", "copy", "-map", "0:v", "-vcodec", "hevc", "output.mkv"},
		},
		{
			name:     "Add Stereo Audio Track",
			inFname:  "input.mkv",
			outFname: "output.mkv",
			inCmd:    []string{"-i", "ENCODARR_INPUT_FILE", "-map", "0:v", "-map", "0:s?", "-map", "0:a", "-map", "0:a", "-c:v", "copy", "-c:s", "copy", "-c:a:1", "copy", "-c:a:0", "aac", "-filter:a:0", "pan=stereo|FL=0.5*FC+0.707*FL+0.707*BL+0.5*LFE|FR=0.5*FC+0.707*FR+0.707*BR+0.5*LFE"},
			expected: []string{"-i", "input.mkv", "-map", "0:v", "-map", "0:s?", "-map", "0:a", "-map", "0:a", "-c:v", "copy", "-c:s", "copy", "-c:a:1", "copy", "-c:a:0", "aac", "-filter:a:0", "pan=stereo|FL=0.5*FC+0.707*FL+0.707*BL+0.5*LFE|FR=0.5*FC+0.707*FR+0.707*BR+0.5*LFE", "output.mkv"},
		},
		{
			name:     "Encode to HEVC and Add Stereo Audio Track",
			inFname:  "input.mkv",
			outFname: "output.mkv",
			inCmd:    []string{"-i", "ENCODARR_INPUT_FILE", "-map", "0:v", "-map", "0:s?", "-map", "0:a", "-map", "0:a", "-c:v", "hevc", "-c:s", "copy", "-c:a:1", "copy", "-c:a:0", "aac", "-filter:a:0", "pan=stereo|FL=0.5*FC+0.707*FL+0.707*BL+0.5*LFE|FR=0.5*FC+0.707*FR+0.707*BR+0.5*LFE"},
			expected: []string{"-i", "input.mkv", "-map", "0:v", "-map", "0:s?", "-map", "0:a", "-map", "0:a", "-c:v", "hevc", "-c:s", "copy", "-c:a:1", "copy", "-c:a:0", "aac", "-filter:a:0", "pan=stereo|FL=0.5*FC+0.707*FL+0.707*BL+0.5*LFE|FR=0.5*FC+0.707*FR+0.707*BR+0.5*LFE", "output.mkv"},
		},
		{
			name:     "Encode to HEVC using Hardware",
			inFname:  "input.mkv",
			outFname: "output.mkv",
			inCmd:    []string{"-hwaccel_device", "/dev/dri/renderD128", "-i", "ENCODARR_INPUT_FILE", "-map", "0:s?", "-map", "0:a", "-c", "copy", "-map", "0:v", "-vcodec", "hevc"},
			expected: []string{"-hwaccel_device", "/dev/dri/renderD128", "-i", "input.mkv", "-map", "0:s?", "-map", "0:a", "-c", "copy", "-map", "0:v", "-vcodec", "hevc", "output.mkv"},
		},
		{
			name:     "All False Params",
			inFname:  "input.mkv",
			outFname: "output.mkv",
			inCmd:    []string{},
			expected: nil,
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.name)

		t.Run(testname, func(t *testing.T) {
			ans := parseFFmpegCmd(tt.inFname, tt.outFname, tt.inCmd)

			if !reflect.DeepEqual(ans, tt.expected) {
				t.Errorf("got %v, expected %v", ans, tt.expected)
			}
		})
	}
}
