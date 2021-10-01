package commanddecider

import (
	"encoding/json"
	"fmt"

	"github.com/BrenekH/encodarr/controller"
)

// codecParams is a map which correlates the TargetVideoCodec settings to the actual parameter to pass to FFMpeg
var codecParams map[string]string = map[string]string{"HEVC": "hevc", "AVC": "libx264", "VP9": "libvpx-vp9"}

// New returns a new CmdDecider.
func New(logger controller.Logger) CmdDecider {
	return CmdDecider{logger: logger}
}

// CmdDecider satisfies the library.CommandDecider interface.
type CmdDecider struct {
	logger controller.Logger
}

// DefaultSettings returns the default settings string.
func (c *CmdDecider) DefaultSettings() string {
	return `{"target_video_codec": "HEVC", "create_stereo_audio": true, "skip_hdr": true, "use_hardware": false, "hardware_codec": "", "hw_device": ""}`
}

// Decide uses the file metadata and settings to decide on a command to run, if any is required.
func (c *CmdDecider) Decide(m controller.FileMetadata, sSettings string) ([]string, error) {
	settings := CmdDeciderSettings{}
	err := json.Unmarshal([]byte(sSettings), &settings)
	if err != nil {
		c.logger.Error(err.Error())
		return []string{}, err
	}

	stereoAudioTrackExists := true
	if settings.CreateStereoAudio {
		for _, v := range m.AudioTracks {
			// stereoAudioTrackExists is reassigned in every loop so that if there are no audio tracks,
			// the decider thinks that there is already a stereo audio track. This still works because
			// if the loop does find a stereo audio track, it breaks instead of running again.
			stereoAudioTrackExists = false

			if v.Channels == 2 {
				stereoAudioTrackExists = true
				break
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

	cmd := genFFmpegCmd(!stereoAudioTrackExists, !alreadyTargetVideoCodec, ffmpegCodecParam, settings.UseHardware, settings.HWDevice)

	return cmd, nil
}

// CmdDeciderSettings defines the structure to unmarshal the settings string into.
type CmdDeciderSettings struct {
	TargetVideoCodec  string `json:"target_video_codec"`
	CreateStereoAudio bool   `json:"create_stereo_audio"`
	SkipHDR           bool   `json:"skip_hdr"`
	UseHardware       bool   `json:"use_hardware"`
	HardwareCodec     string `json:"hardware_codec"`
	HWDevice          string `json:"hw_device"`
}

// genFFmpegCmd creates the correct ffmpeg arguments for the input/output filenames and the job parameters.
func genFFmpegCmd(stereo, encode bool, codec string, useHW bool, hwDevice string) []string {
	var s []string

	if stereo && encode {
		s = []string{"-i", "ENCODARR_INPUT_FILE", "-map", "0:v", "-map", "0:s?", "-map", "0:a", "-map", "0:a", "-c:v", codec, "-c:s", "copy", "-c:a:1", "copy", "-c:a:0", "aac", "-filter:a:0", "pan=stereo|FL=0.5*FC+0.707*FL+0.707*BL+0.5*LFE|FR=0.5*FC+0.707*FR+0.707*BR+0.5*LFE"}
	} else if stereo {
		s = []string{"-i", "ENCODARR_INPUT_FILE", "-map", "0:v", "-map", "0:s?", "-map", "0:a", "-map", "0:a", "-c:v", "copy", "-c:s", "copy", "-c:a:1", "copy", "-c:a:0", "aac", "-filter:a:0", "pan=stereo|FL=0.5*FC+0.707*FL+0.707*BL+0.5*LFE|FR=0.5*FC+0.707*FR+0.707*BR+0.5*LFE"}
	} else if encode {
		s = []string{"-i", "ENCODARR_INPUT_FILE", "-map", "0:s?", "-map", "0:a", "-c", "copy", "-map", "0:v", "-vcodec", codec}
	}

	if hwDevice != "" && useHW {
		s = append([]string{"-hwaccel_device", hwDevice}, s...)
	}

	return s
}
