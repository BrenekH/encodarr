package mediainfo

import "github.com/BrenekH/encodarr/controller"

type MetadataReader struct {
	logger controller.Logger
}

func (m *MetadataReader) Read(path string) (fm controller.FileMetadata) {
	m.logger.Critical("Not implemented")
	// TODO: Implement
	return
}
