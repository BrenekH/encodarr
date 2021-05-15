package mediainfo

import "github.com/BrenekH/encodarr/controller"

func NewMetadataReader(logger controller.Logger) MetadataReader {
	return MetadataReader{logger: logger}
}

type MetadataReader struct {
	logger controller.Logger
}

func (m *MetadataReader) Read(path string) (fm controller.FileMetadata) {
	m.logger.Critical("Not implemented")
	// TODO: Implement
	return
}
