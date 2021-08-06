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
			if generalDuration, err = strconv.ParseFloat(v.Duration, 32); err != nil {
				m.logger.Debug("error while parsing general duration (%v) for %v: %v", path, v.Duration, err)
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

			if vidTrack.Index, err = strconv.Atoi(v.StreamOrder); err != nil {
				m.logger.Debug("error while converting vidTrack.Index (StreamOrder) for %v: %v", path, err)
				return controller.FileMetadata{}, err
			}

			if vidTrack.Width, err = strconv.Atoi(v.Width); err != nil {
				m.logger.Debug("error while converting vidTrack.Width for %v: %v", path, err)
				return controller.FileMetadata{}, err
			}

			if vidTrack.Height, err = strconv.Atoi(v.Height); err != nil {
				m.logger.Debug("error while converting vidTrack.Height for %v: %v", path, err)
				return controller.FileMetadata{}, err
			}

			vidTracks = append(vidTracks, vidTrack)
		case "Audio":
			audioTrack := controller.AudioTrack{}

			if audioTrack.Index, err = strconv.Atoi(v.StreamOrder); err != nil {
				m.logger.Debug("error while converting audioTrack.Index (StreamOrder) for %v: %v", path, err)
				return controller.FileMetadata{}, err
			}

			if audioTrack.Channels, err = strconv.Atoi(v.Channels); err != nil {
				m.logger.Debug("error while converting audioTrack.Channels for %v: %v", path, err)
				return controller.FileMetadata{}, err
			}

			audioTracks = append(audioTracks, audioTrack)
		case "Text":
			textTrack := controller.SubtitleTrack{}

			if textTrack.Index, err = strconv.Atoi(v.StreamOrder); err != nil {
				if textTrack.Index, err = strconv.Atoi(v.UniqueID); err != nil {
					m.logger.Warn("error while converting textTrack.Index (StreamOrder, UniqueID) for %v: %v", path, err)
					continue
				}
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
