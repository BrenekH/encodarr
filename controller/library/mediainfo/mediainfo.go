package mediainfo

import (
	"encoding/json"
	"strconv"

	"github.com/BrenekH/encodarr/controller"
)

func NewMetadataReader(logger controller.Logger) MetadataReader {
	return MetadataReader{
		logger: logger,
		cmdr:   ExecCommander{},
	}
}

type MetadataReader struct {
	logger controller.Logger

	cmdr Commander
}

func (m *MetadataReader) Read(path string) (controller.FileMetadata, error) {
	cmd := m.cmdr.Command("mediainfo", "--Output=JSON", "--Full", path)
	b, err := cmd.Output()
	if err != nil {
		return controller.FileMetadata{}, err
	}

	mi := mediaInfo{}
	err = json.Unmarshal(b, &mi)
	if err != nil {
		return controller.FileMetadata{}, err
	}

	var generalDuration float64
	vidTracks := make([]controller.VideoTrack, 0)
	audioTracks := make([]controller.AudioTrack, 0)
	subtitleTracks := make([]controller.SubtitleTrack, 0)

	for _, v := range mi.Media.Tracks {
		switch v.Type {
		case "General":
			generalDuration, err = strconv.ParseFloat(v.Duration, 32)
			if err != nil {
				m.logger.Debug("error while parsing general duration (%v): %v", v.Duration, err)
				return controller.FileMetadata{}, err
			}
		case "Video":
			vidTrack := controller.VideoTrack{}

			switch v.Format {
			case "AVC":
				vidTrack.Codec = "AVC"
			case "HEVC":
				vidTrack.Codec = "HEVC"
			case "VP9":
				vidTrack.Codec = "VP9"
			case "AV1":
				vidTrack.Codec = "AV1"
			default:
				vidTrack.Codec = ""
			}

			vidTrack.ColorPrimaries = v.ColourPrimaries

			vidTrack.Index, err = strconv.Atoi(v.StreamOrder)
			if err != nil {
				m.logger.Debug("error while converting vidTrack.Index (StreamOrder: %v): %v", v.StreamOrder, err)
				return controller.FileMetadata{}, err
			}

			vidTrack.Width, err = strconv.Atoi(v.Width)
			if err != nil {
				m.logger.Debug("error while converting vidTrack.Width (%v): %v", v.Width, err)
				return controller.FileMetadata{}, err
			}

			vidTrack.Height, err = strconv.Atoi(v.Height)
			if err != nil {
				m.logger.Debug("error while converting vidTrack.Height (%v): %v", v.Height, err)
				return controller.FileMetadata{}, err
			}

			vidTracks = append(vidTracks, vidTrack)
		case "Audio":
			audioTrack := controller.AudioTrack{}

			audioTrack.Index, err = strconv.Atoi(v.StreamOrder)
			if err != nil {
				m.logger.Debug("error while converting audioTrack.Index (StreamOrder: %v): %v", v.StreamOrder, err)
				return controller.FileMetadata{}, err
			}

			audioTrack.Channels, err = strconv.Atoi(v.Channels)
			if err != nil {
				m.logger.Debug("error while converting audioTrack.Channels (%v): %v", v.Channels, err)
				return controller.FileMetadata{}, err
			}

			audioTracks = append(audioTracks, audioTrack)
		case "Text":
			textTrack := controller.SubtitleTrack{}

			textTrack.Index, err = strconv.Atoi(v.StreamOrder)
			if err != nil {
				m.logger.Debug("error while converting textTrack.Index (StreamOrder: %v): %v", v.StreamOrder, err)
				return controller.FileMetadata{}, err
			}

			textTrack.Language = v.Language

			subtitleTracks = append(subtitleTracks, textTrack)
		case "Menu":
		default:
		}

	}

	return controller.FileMetadata{
		General: controller.General{
			Duration: float32(generalDuration),
		},
		VideoTracks:    vidTracks,
		AudioTracks:    audioTracks,
		SubtitleTracks: subtitleTracks,
	}, nil
}
