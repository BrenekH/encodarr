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
			params:   jobParameters{Encode: true, Codec: "hevc", Stereo: false},
			expected: []string{"-map", "0:s?", "-map", "0:a", "-c", "copy", "-map", "0:v", "-vcodec", "hevc"},
		},
		{
			name:     "Add Stereo Audio Track",
			params:   jobParameters{Stereo: true, Encode: false, Codec: ""},
			expected: []string{"-map", "0:v", "-map", "0:s?", "-map", "0:a", "-map", "0:a", "-c:v", "copy", "-c:s", "copy", "-c:a:1", "copy", "-c:a:0", "aac", "-filter:a:0", "pan=stereo|FL=0.5*FC+0.707*FL+0.707*BL+0.5*LFE|FR=0.5*FC+0.707*FR+0.707*BR+0.5*LFE"},
		},
		{
			name:     "Encode to HEVC and Add Stereo Audio Track",
			params:   jobParameters{Encode: true, Codec: "hevc", Stereo: true},
			expected: []string{"-map", "0:v", "-map", "0:s?", "-map", "0:a", "-map", "0:a", "-c:v", "hevc", "-c:s", "copy", "-c:a:1", "copy", "-c:a:0", "aac", "-filter:a:0", "pan=stereo|FL=0.5*FC+0.707*FL+0.707*BL+0.5*LFE|FR=0.5*FC+0.707*FR+0.707*BR+0.5*LFE"},
		},
		{
			name:     "All False Params",
			params:   jobParameters{Encode: false, Codec: "", Stereo: false},
			expected: nil,
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.name)

		t.Run(testname, func(t *testing.T) {
			ans := genFFmpegCmd(tt.params.Stereo, tt.params.Encode, tt.params.Codec)

			if !reflect.DeepEqual(ans, tt.expected) {
				t.Errorf("got %v, expected %v", ans, tt.expected)
			}
		})
	}
}

type jobParameters struct {
	Stereo bool
	Encode bool
	Codec  string
}
