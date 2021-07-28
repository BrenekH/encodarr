package http

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGenFFmpegCmd(t *testing.T) {
	tests := []struct {
		name     string
		inFname  string
		outFname string
		params   jobParameters
		expected []string
	}{
		{
			name:     "Encode to HEVC",
			inFname:  "input.mkv",
			outFname: "output.mkv",
			params:   jobParameters{Encode: true, Codec: "hevc", Stereo: false, HWDevice: ""},
			expected: []string{"-i", "input.mkv", "-map", "0:s?", "-map", "0:a", "-c", "copy", "-map", "0:v", "-vcodec", "hevc", "output.mkv"},
		},
		{
			name:     "Add Stereo Audio Track",
			inFname:  "input.mkv",
			outFname: "output.mkv",
			params:   jobParameters{Stereo: true, Encode: false, Codec: "", HWDevice: ""},
			expected: []string{"-i", "input.mkv", "-map", "0:v", "-map", "0:s?", "-map", "0:a", "-map", "0:a", "-c:v", "copy", "-c:s", "copy", "-c:a:1", "copy", "-c:a:0", "aac", "-filter:a:0", "pan=stereo|FL=0.5*FC+0.707*FL+0.707*BL+0.5*LFE|FR=0.5*FC+0.707*FR+0.707*BR+0.5*LFE", "output.mkv"},
		},
		{
			name:     "Encode to HEVC and Add Stereo Audio Track",
			inFname:  "input.mkv",
			outFname: "output.mkv",
			params:   jobParameters{Encode: true, Codec: "hevc", Stereo: true, HWDevice: ""},
			expected: []string{"-i", "input.mkv", "-map", "0:v", "-map", "0:s?", "-map", "0:a", "-map", "0:a", "-c:v", "hevc", "-c:s", "copy", "-c:a:1", "copy", "-c:a:0", "aac", "-filter:a:0", "pan=stereo|FL=0.5*FC+0.707*FL+0.707*BL+0.5*LFE|FR=0.5*FC+0.707*FR+0.707*BR+0.5*LFE", "output.mkv"},
		},
		{
			name:     "All False Params",
			inFname:  "input.mkv",
			outFname: "output.mkv",
			params:   jobParameters{Encode: false, Codec: "", Stereo: false, HWDevice: ""},
			expected: nil,
		},
		{
			name:     "Use hardware for encoding",
			inFname:  "input.mkv",
			outFname: "output.mkv",
			params:   jobParameters{Encode: true, Codec: "hevc_vaapi", Stereo: false, HWDevice: "/dev/gpu1"},
			expected: []string{"-hwaccel_device", "/dev/gpu1", "-i", "input.mkv", "-map", "0:s?", "-map", "0:a", "-c", "copy", "-map", "0:v", "-vcodec", "hevc_vaapi", "output.mkv"},
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.name)

		t.Run(testname, func(t *testing.T) {
			ans := genFFmpegCmd(tt.inFname, tt.outFname, tt.params)

			if !reflect.DeepEqual(ans, tt.expected) {
				t.Errorf("got %v, expected %v", ans, tt.expected)
			}
		})
	}
}
