package commanddecider

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/BrenekH/encodarr/controller"
)

func TestDefaultSettingsUnmarshals(t *testing.T) {
	cmdDecider := CmdDecider{}
	defaultSettingsString := cmdDecider.DefaultSettings()
	cmdDeciderSettings := CmdDeciderSettings{}

	err := json.Unmarshal([]byte(defaultSettingsString), &cmdDeciderSettings)

	if err != nil {
		t.Errorf("unexpected error while unmarshaling CmdDecider.DefaultSettings() into CmdDeciderSettings: '%v'", err)
	}
}

func TestCmdDeciderDecide(t *testing.T) {
	// TODO: Other test cases
	//   - There are no video tracks

	tests := []struct {
		name        string
		metadata    controller.FileMetadata
		settings    string
		errExpected bool
		expected    []string
	}{
		{
			name: "Encode to HEVC",
			metadata: controller.FileMetadata{
				VideoTracks: []controller.VideoTrack{
					{
						Codec: "AVC",
					},
				},
			},
			settings: `{"target_video_codec": "HEVC", "create_stereo_audio": true, "skip_hdr": true, "use_hardware": false, "hardware_codec": "", "hw_device": ""}`,
			expected: []string{"-i", "ENCODARR_INPUT_FILE", "-map", "0:s?", "-map", "0:a", "-c", "copy", "-map", "0:v", "-vcodec", "hevc"},
		},
		{
			name: "Add Stereo Audio Track",
			metadata: controller.FileMetadata{
				VideoTracks: []controller.VideoTrack{
					{
						Codec: "HEVC",
					},
				},
				AudioTracks: []controller.AudioTrack{
					{
						Channels: 6,
					},
				},
			},
			settings: `{"target_video_codec": "HEVC", "create_stereo_audio": true, "skip_hdr": true, "use_hardware": false, "hardware_codec": "", "hw_device": ""}`,
			expected: []string{"-i", "ENCODARR_INPUT_FILE", "-map", "0:v", "-map", "0:s?", "-map", "0:a", "-map", "0:a", "-c:v", "copy", "-c:s", "copy", "-c:a:1", "copy", "-c:a:0", "aac", "-filter:a:0", "pan=stereo|FL=0.5*FC+0.707*FL+0.707*BL+0.5*LFE|FR=0.5*FC+0.707*FR+0.707*BR+0.5*LFE"},
		},
		{
			name: "Encode to HEVC and Add Stereo Audio Track",
			metadata: controller.FileMetadata{
				VideoTracks: []controller.VideoTrack{
					{
						Codec: "AVC",
					},
				},
				AudioTracks: []controller.AudioTrack{
					{
						Channels: 6,
					},
				},
			},
			settings: `{"target_video_codec": "HEVC", "create_stereo_audio": true, "skip_hdr": true, "use_hardware": false, "hardware_codec": "", "hw_device": ""}`,
			expected: []string{"-i", "ENCODARR_INPUT_FILE", "-map", "0:v", "-map", "0:s?", "-map", "0:a", "-map", "0:a", "-c:v", "hevc", "-c:s", "copy", "-c:a:1", "copy", "-c:a:0", "aac", "-filter:a:0", "pan=stereo|FL=0.5*FC+0.707*FL+0.707*BL+0.5*LFE|FR=0.5*FC+0.707*FR+0.707*BR+0.5*LFE"},
		},
		{
			name: "Encode to HEVC using hardware",
			metadata: controller.FileMetadata{
				VideoTracks: []controller.VideoTrack{
					{
						Codec: "AVC",
					},
				},
			},
			settings: `{"target_video_codec": "HEVC", "create_stereo_audio": true, "skip_hdr": true, "use_hardware": true, "hardware_codec": "hevc_vaapi", "hw_device": "/dev/dri/renderD128"}`,
			expected: []string{"-hwaccel_device", "/dev/dri/renderD128", "-i", "ENCODARR_INPUT_FILE", "-map", "0:s?", "-map", "0:a", "-c", "copy", "-map", "0:v", "-vcodec", "hevc_vaapi"},
		},
		{
			name: "All False Params",
			metadata: controller.FileMetadata{
				VideoTracks: []controller.VideoTrack{
					{
						Codec: "HEVC",
					},
				},
			},
			settings:    `{"target_video_codec": "HEVC", "create_stereo_audio": true, "skip_hdr": true, "use_hardware": false, "hardware_codec": "", "hw_device": ""}`,
			errExpected: true,
			expected:    []string{},
		},
		{
			name: "create_stereo_audio is true and a stereo track already exists",
			metadata: controller.FileMetadata{
				VideoTracks: []controller.VideoTrack{
					{
						Codec: "HEVC",
					},
				},
				AudioTracks: []controller.AudioTrack{
					{
						Channels: 2,
					},
				},
			},
			settings:    `{"target_video_codec": "HEVC", "create_stereo_audio": true, "skip_hdr": true, "use_hardware": false, "hardware_codec": "", "hw_device": ""}`,
			errExpected: true,
			expected:    []string{},
		},
		{
			name: "target_video_codec is one that is not a part of the codecParams map",
			metadata: controller.FileMetadata{
				VideoTracks: []controller.VideoTrack{
					{
						Codec: "AVC",
					},
				},
			},
			settings:    `{"target_video_codec": "CODEC", "create_stereo_audio": false, "skip_hdr": true, "use_hardware": false, "hardware_codec": "", "hw_device": ""}`,
			errExpected: true,
			expected:    []string{},
		},
		{
			name:        "Settings string is invalid",
			metadata:    controller.FileMetadata{},
			settings:    `target_video_codec": "CODEC", "create_stereo_audio": false, "skip_hdr": true, "use_hardware": false, "hardware_codec": "", "hw_device": ""}`,
			errExpected: true,
			expected:    []string{},
		},
		{
			name: "Add Stereo Audio Track (No Video Tracks)",
			metadata: controller.FileMetadata{
				AudioTracks: []controller.AudioTrack{
					{
						Channels: 6,
					},
				},
			},
			settings: `{"target_video_codec": "HEVC", "create_stereo_audio": true, "skip_hdr": true, "use_hardware": false, "hardware_codec": "", "hw_device": ""}`,
			expected: []string{"-i", "ENCODARR_INPUT_FILE", "-map", "0:v", "-map", "0:s?", "-map", "0:a", "-map", "0:a", "-c:v", "copy", "-c:s", "copy", "-c:a:1", "copy", "-c:a:0", "aac", "-filter:a:0", "pan=stereo|FL=0.5*FC+0.707*FL+0.707*BL+0.5*LFE|FR=0.5*FC+0.707*FR+0.707*BR+0.5*LFE"},
		},
		{
			name: "No Video Track and 2-Channel Audio Track",
			metadata: controller.FileMetadata{
				AudioTracks: []controller.AudioTrack{
					{
						Channels: 2,
					},
				},
			},
			settings:    `{"target_video_codec": "HEVC", "create_stereo_audio": true, "skip_hdr": true, "use_hardware": false, "hardware_codec": "", "hw_device": ""}`,
			errExpected: true,
			expected:    []string{},
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.name)
		mLogger := mockLogger{}
		cmdDecider := CmdDecider{logger: &mLogger}

		t.Run(testname, func(t *testing.T) {
			ans, err := cmdDecider.Decide(tt.metadata, tt.settings)

			if err != nil && !tt.errExpected {
				t.Errorf("error was expected to be nil, but instead it was '%v'", err)
			}

			if !reflect.DeepEqual(ans, tt.expected) {
				t.Errorf("got %v, expected %v", ans, tt.expected)
			}
		})
	}
}

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

type mockLogger struct{}

func (m *mockLogger) Trace(s string, i ...interface{})    {}
func (m *mockLogger) Debug(s string, i ...interface{})    {}
func (m *mockLogger) Info(s string, i ...interface{})     {}
func (m *mockLogger) Warn(s string, i ...interface{})     {}
func (m *mockLogger) Error(s string, i ...interface{})    {}
func (m *mockLogger) Critical(s string, i ...interface{}) {}
