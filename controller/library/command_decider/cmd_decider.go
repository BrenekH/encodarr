package command_decider

import (
	"encoding/json"
	"fmt"

	"github.com/BrenekH/encodarr/controller"
)

// codecParams is a map which correlates the TargetVideoCodec settings to the actual parameter to pass to FFMpeg
var codecParams map[string]string = map[string]string{"HEVC": "hevc", "AVC": "libx264", "VP9": "libvpx-vp9"}

func New(logger controller.Logger) CmdDecider {
	return CmdDecider{logger: logger}
}

type CmdDecider struct {
	logger controller.Logger
}

func (c *CmdDecider) DefaultSettings() string {
	return `{"target_video_codec": "HEVC", "create_stereo_audio": true, "skip_hdr": true, "use_hardware": false, "hardware_codec": "", "hw_device": ""}`
}

func (c *CmdDecider) Decide(m controller.FileMetadata, sSettings string) ([]string, error) {
	settings := CmdDeciderSettings{}
	err := json.Unmarshal([]byte(sSettings), &settings)
	if err != nil {
		c.logger.Error(err.Error())
		return []string{}, err
	}

	stereoAudioTrackExists := true
	if settings.CreateStereoAudio {
		stereoAudioTrackExists = false
		for _, v := range m.AudioTracks {
			if v.Channels == 2 {
				stereoAudioTrackExists = true
			}
		}
	}

	var alreadyTargetVideoCodec bool
	if len(m.VideoTracks) > 0 {
		alreadyTargetVideoCodec = m.VideoTracks[0].Codec == settings.TargetVideoCodec
	} else {
		// Just because there are no video tracks, doesn't mean that the audio can't be adjusted.
		// So tell the system that the video is already the target and move on.
		alreadyTargetVideoCodec = true
	}

	if stereoAudioTrackExists && alreadyTargetVideoCodec {
		return []string{}, fmt.Errorf("file already matches requirements")
	}

	var ffmpegCodecParam string
	if settings.UseHardware {
		ffmpegCodecParam = settings.HardwareCodec
	} else {
		var ok bool
		ffmpegCodecParam, ok = codecParams[settings.TargetVideoCodec]
		if !ok {
			return []string{}, fmt.Errorf("couldn't identify ffmpeg parameter for '%v' target codec", settings.TargetVideoCodec)
		}
	}

	cmd := genFFmpegCmd(!stereoAudioTrackExists, !alreadyTargetVideoCodec, ffmpegCodecParam, settings.HWDevice)

	return cmd, nil
}

type CmdDeciderSettings struct {
	TargetVideoCodec  string `json:"target_video_codec"`
	CreateStereoAudio bool   `json:"create_stereo_audio"`
	SkipHDR           bool   `json:"skip_hdr"`
	UseHardware       bool   `bool:"use_hardware"`
	HardwareCodec     string `json:"hardware_codec"`
	HWDevice          string `json:"hw_device"`
}

// genFFmpegCmd creates the correct ffmpeg arguments for the input/output filenames and the job parameters.
func genFFmpegCmd(stereo, encode bool, codec, hwDevice string) []string {
	var s []string
	if stereo && encode {
		s = []string{"-map", "0:v", "-map", "0:s?", "-map", "0:a", "-map", "0:a", "-c:v", codec, "-c:s", "copy", "-c:a:1", "copy", "-c:a:0", "aac", "-filter:a:0", "pan=stereo|FL=0.5*FC+0.707*FL+0.707*BL+0.5*LFE|FR=0.5*FC+0.707*FR+0.707*BR+0.5*LFE"}
	} else if stereo {
		s = []string{"-map", "0:v", "-map", "0:s?", "-map", "0:a", "-map", "0:a", "-c:v", "copy", "-c:s", "copy", "-c:a:1", "copy", "-c:a:0", "aac", "-filter:a:0", "pan=stereo|FL=0.5*FC+0.707*FL+0.707*BL+0.5*LFE|FR=0.5*FC+0.707*FR+0.707*BR+0.5*LFE"}
	} else if encode {
		s = []string{"-map", "0:s?", "-map", "0:a", "-c", "copy", "-map", "0:v", "-vcodec", codec}
	}
	return s
}
