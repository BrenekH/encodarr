package command_decider

import (
	"fmt"
	"reflect"
	"testing"
)

// TODO: Test CmdDecider.Decide

func TestGenFFmpegCmd(t *testing.T) {
	tests := []struct {
		name     string
		params   jobParameters
		expected []string
	}{
		{
			name:     "Encode to HEVC",
			params:   jobParameters{Encode: true, Codec: "hevc", Stereo: false, UseHW: false, HWDevice: ""},
			expected: []string{"-i", "ENCODARR_INPUT_FILE", "-map", "0:s?", "-map", "0:a", "-c", "copy", "-map", "0:v", "-vcodec", "hevc"},
		},
		{
			name:     "Add Stereo Audio Track",
			params:   jobParameters{Stereo: true, Encode: false, Codec: "", UseHW: false, HWDevice: ""},
			expected: []string{"-i", "ENCODARR_INPUT_FILE", "-map", "0:v", "-map", "0:s?", "-map", "0:a", "-map", "0:a", "-c:v", "copy", "-c:s", "copy", "-c:a:1", "copy", "-c:a:0", "aac", "-filter:a:0", "pan=stereo|FL=0.5*FC+0.707*FL+0.707*BL+0.5*LFE|FR=0.5*FC+0.707*FR+0.707*BR+0.5*LFE"},
		},
		{
			name:     "Encode to HEVC and Add Stereo Audio Track",
			params:   jobParameters{Encode: true, Codec: "hevc", Stereo: true, UseHW: false, HWDevice: ""},
			expected: []string{"-i", "ENCODARR_INPUT_FILE", "-map", "0:v", "-map", "0:s?", "-map", "0:a", "-map", "0:a", "-c:v", "hevc", "-c:s", "copy", "-c:a:1", "copy", "-c:a:0", "aac", "-filter:a:0", "pan=stereo|FL=0.5*FC+0.707*FL+0.707*BL+0.5*LFE|FR=0.5*FC+0.707*FR+0.707*BR+0.5*LFE"},
		},
		{
			name:     "Encode to HEVC using hardware",
			params:   jobParameters{Encode: true, Codec: "hevc_vaapi", Stereo: false, UseHW: true, HWDevice: "/dev/dri/renderD128"},
			expected: []string{"-hwaccel_device", "/dev/dri/renderD128", "-i", "ENCODARR_INPUT_FILE", "-map", "0:s?", "-map", "0:a", "-c", "copy", "-map", "0:v", "-vcodec", "hevc_vaapi"},
		},
		{
			name:     "Don't add hwaccel_device if not using hardware encoding",
			params:   jobParameters{Encode: true, Codec: "hevc_vaapi", Stereo: false, UseHW: false, HWDevice: "/dev/dri/renderD128"},
			expected: []string{"-i", "ENCODARR_INPUT_FILE", "-map", "0:s?", "-map", "0:a", "-c", "copy", "-map", "0:v", "-vcodec", "hevc_vaapi"},
		},
		{
			name:     "All False Params",
			params:   jobParameters{Encode: false, Codec: "", Stereo: false, UseHW: false, HWDevice: ""},
			expected: nil,
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.name)

		t.Run(testname, func(t *testing.T) {
			ans := genFFmpegCmd(tt.params.Stereo, tt.params.Encode, tt.params.Codec, tt.params.UseHW, tt.params.HWDevice)

			if !reflect.DeepEqual(ans, tt.expected) {
				t.Errorf("got %v, expected %v", ans, tt.expected)
			}
		})
	}
}

type jobParameters struct {
	Stereo   bool
	Encode   bool
	Codec    string
	UseHW    bool
	HWDevice string
}
