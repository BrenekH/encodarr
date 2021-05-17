package mediainfo

import (
	"encoding/json"

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

	return controller.FileMetadata{
		General: controller.General{
			Duration: 0,
		},
		VideoTracks:    []controller.VideoTrack{},
		AudioTracks:    []controller.AudioTrack{},
		SubtitleTracks: []controller.SubtitleTrack{},
	}, nil
}
