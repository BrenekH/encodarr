package library

import (
	"io/fs"
	"os"
	"testing"
	"time"

	"github.com/BrenekH/encodarr/controller"
)

func TestNewCacheSetsInternalFields(t *testing.T) {
	//? I'm not sure that using a uniqueId is the best way to validate that the internal fields are set properly.
	//? Maybe instead it should check the object locations or something similar to validate.

	m := mockMetadataReader{uniqueId: "testId1"}
	f := mockFileCacheDataStorer{uniqueId: "testId2"}
	l := mockLogger{uniqueId: "testId3"}

	receivedStruct := NewCache(&m, &f, &l)

	if (*(receivedStruct.metadataReader.(*mockMetadataReader))).uniqueId != m.uniqueId {
		t.Errorf("MetadataReader inside Cache struct is not the same one that was passed to NewCache")
	}

	if (*(receivedStruct.ds.(*mockFileCacheDataStorer))).uniqueId != f.uniqueId {
		t.Errorf("FileCacheDataStorer inside Cache struct is not the same one that was passed to NewCache")
	}

	if (*(receivedStruct.logger.(*mockLogger))).uniqueId != l.uniqueId {
		t.Errorf("Logger inside Cache struct is not the same one that was passed to NewCache")
	}
}

// TODO: Test Cache
//   = Read - Reading from metadataReader when the modtimes are different. Returning stored metadata when modtimes are same.

func TestCacheReadDifferentModtimes(t *testing.T) {
	m := mockMetadataReader{}
	f := mockFileCacheDataStorer{}
	l := mockLogger{}

	cache := NewCache(&m, &f, &l)
	s := mockStater{}
	cache.stater = &s

	// TODO: Check that metadata reader's Read function is called.
}

func TestCacheReadSameModtimes(t *testing.T) {
	m := mockMetadataReader{}
	f := mockFileCacheDataStorer{}
	l := mockLogger{}

	cache := NewCache(&m, &f, &l)
	s := mockStater{}
	cache.stater = &s

	// TODO: Check that data storer's Read function is called.
}

type mockMetadataReader struct {
	uniqueId string
}

func (m *mockMetadataReader) Read(path string) (controller.FileMetadata, error) {
	return controller.FileMetadata{}, nil
}

type mockFileCacheDataStorer struct {
	uniqueId string
}

func (m *mockFileCacheDataStorer) Modtime(path string) (time.Time, error) {
	return time.Now(), nil
}

func (m *mockFileCacheDataStorer) Metadata(path string) (controller.FileMetadata, error) {
	return controller.FileMetadata{}, nil
}

func (m *mockFileCacheDataStorer) SaveModtime(path string, t time.Time) error {
	return nil
}

func (m *mockFileCacheDataStorer) SaveMetadata(path string, f controller.FileMetadata) error {
	return nil
}

type mockLogger struct {
	uniqueId string
}

func (m *mockLogger) Trace(s string, i ...interface{})    {}
func (m *mockLogger) Debug(s string, i ...interface{})    {}
func (m *mockLogger) Info(s string, i ...interface{})     {}
func (m *mockLogger) Warn(s string, i ...interface{})     {}
func (m *mockLogger) Error(s string, i ...interface{})    {}
func (m *mockLogger) Critical(s string, i ...interface{}) {}

type mockStater struct {
}

func (m *mockStater) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}
