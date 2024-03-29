package library

import (
	"database/sql"
	"io/fs"
	"os"
	"time"

	"github.com/BrenekH/encodarr/controller"
)

// NewCache returns a new Cache.
func NewCache(m MetadataReader, f controller.FileCacheDataStorer, l controller.Logger) Cache {
	return Cache{
		metadataReader: m,
		ds:             f,
		logger:         l,
		stater:         osStater{},
	}
}

// Cache sits in front of a MetadataReader and only calls it for
// a Read call when the file has updated(based on the modtime)
type Cache struct {
	metadataReader MetadataReader
	ds             controller.FileCacheDataStorer
	logger         controller.Logger
	stater         stater
}

// Read uses the data storer and file.Stat to determine whether or not to call the MetadataReader or return from the cache.
func (c *Cache) Read(path string) (controller.FileMetadata, error) {
	fileInfo, err := c.stater.Stat(path)
	if err != nil {
		c.logger.Error("Failed to stat %v, disabling caching for this call: %v", path, err)
		return c.metadataReader.Read(path)
	}

	storedModtime, err := c.ds.Modtime(path)
	if err != nil {
		if err != sql.ErrNoRows {
			c.logger.Error("Failed to read stored modtime for %v, disabling caching for this call: %v", path, err)
			return c.metadataReader.Read(path)
		}
		storedModtime = time.Unix(0, 0)
	}

	// We have to set the mod times to UTC because the db returns a different time zone format than os.Stat()
	if fileInfo.ModTime().UTC() == storedModtime.UTC() {
		storedMetadata, err := c.ds.Metadata(path)
		if err != nil {
			c.logger.Error("Failed to read stored metadata for %v, disabling caching for this call: %v", path, err)
			return c.metadataReader.Read(path)
		}

		return storedMetadata, nil
	}

	newMetadata, err := c.metadataReader.Read(path)
	if err == nil {
		err = c.ds.SaveMetadata(path, newMetadata)
		if err != nil {
			c.logger.Error("Failed to save new metadata for %v: %v", path, err)
		}

		err = c.ds.SaveModtime(path, fileInfo.ModTime())
		if err != nil {
			c.logger.Error("Failed to save new modtime for %v: %v", path, err)
		}
	}

	return newMetadata, err
}

type osStater struct{}

func (o osStater) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}
