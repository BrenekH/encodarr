package library

import "github.com/BrenekH/encodarr/controller"

func NewCache(m MetadataReader) Cache {
	return Cache{metadataReader: m}
}

// Cache sits in front of a MetadataReader and only calls it for
// a Read call when the file has updated(via modtime)
type Cache struct {
	metadataReader MetadataReader
}

func (c *Cache) Read(path string) (controller.FileMetadata, error) {
	// TODO: Implement caching
	return c.metadataReader.Read(path)
}
