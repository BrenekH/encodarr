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

func (c *CmdDecider) Decide(m controller.FileMetadata, bSettings []byte) (bool, []string) { //? Maybe an extra error return value would be good?
	settings := CmdDeciderSettings{}
	err := json.Unmarshal(bSettings, &settings)
	if err != nil {
		c.logger.Error(err.Error())
		return false, []string{"Error parsing json settings", err.Error()}
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

	alreadyTargetVideoCodec := m.VideoTracks[0].Codec == settings.TargetVideoCodec // TODO: Protect against empty slice

	if stereoAudioTrackExists && alreadyTargetVideoCodec {
		return false, []string{"File already matches requirements"}
	}

	ffmpegCodecParam, ok := codecParams[settings.TargetVideoCodec]
	if !ok {
		return false, []string{fmt.Sprintf("Couldn't identify ffmpeg parameter for '%v' target codec", settings.TargetVideoCodec)}
	}

	cmd := genFFmpegCmd(!stereoAudioTrackExists, !alreadyTargetVideoCodec, ffmpegCodecParam)

	return true, cmd
}

type CmdDeciderSettings struct {
	TargetVideoCodec  string
	CreateStereoAudio bool
	SkipHDR           bool
}

// genFFmpegCmd creates the correct ffmpeg arguments for the input/output filenames and the job parameters.
func genFFmpegCmd(stereo, encode bool, codec string) []string {
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