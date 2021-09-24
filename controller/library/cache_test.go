package library

import (
	"io/fs"
	"os"
	"testing"
	"time"

	"github.com/BrenekH/encodarr/controller"
)

// TODO: Test Cache.Read when stater and data storer error.

func TestNewCacheSetsInternalFields(t *testing.T) {
	//? I'm not sure that using a uniqueId is the best way to validate that the internal fields are set properly.
	//? Maybe instead it should check the object locations or something similar to validate.

	m := mockMetadataReader{uniqueId: "testId1"}
	f := mockFileCacheDataStorer{uniqueId: "testId2"}
	l := mockLogger{uniqueId: "testId3"}

	receivedStruct := NewCache(&m, &f, &l)

	if (*(receivedStruct.metadataReader.(*mockMetadataReader))).uniqueId != m.uniqueId {
		t.Error("MetadataReader inside Cache struct is not the same one that was passed to NewCache")
	}

	if (*(receivedStruct.ds.(*mockFileCacheDataStorer))).uniqueId != f.uniqueId {
		t.Error("FileCacheDataStorer inside Cache struct is not the same one that was passed to NewCache")
	}

	if (*(receivedStruct.logger.(*mockLogger))).uniqueId != l.uniqueId {
		t.Error("Logger inside Cache struct is not the same one that was passed to NewCache")
	}
}

func TestCacheReadDifferentModtimes(t *testing.T) {
	m := mockMetadataReader{}
	f := mockFileCacheDataStorer{modtimeReturnData: time.Unix(10000, 100)}
	l := mockLogger{}

	cache := NewCache(&m, &f, &l)
	s := mockStater{statReturnValue: &fileInfo{modtime: time.Unix(20000, 200)}}
	cache.stater = &s

	cache.Read("test")

	if !m.readCalled {
		t.Error("MetadataReader.Read() was not called even though it should have been")
	}

	if f.metadataCalled {
		t.Error("FileCacheDataStorer.Metadata() was called even though it should not have been")
	}
}

func TestCacheReadSameModtimes(t *testing.T) {
	mtime := time.Unix(10000, 100)

	m := mockMetadataReader{}
	f := mockFileCacheDataStorer{modtimeReturnData: mtime}
	l := mockLogger{}

	cache := NewCache(&m, &f, &l)
	s := mockStater{statReturnValue: &fileInfo{modtime: mtime}}
	cache.stater = &s

	cache.Read("test")

	if !f.metadataCalled {
		t.Error("FileCacheDataStorer.Metdata() was not called even though it should have been")
	}

	if m.readCalled {
		t.Error("MetadataReader.Read() was called even though it should not have been")
	}
}

type mockMetadataReader struct {
	uniqueId       string
	readReturnData controller.FileMetadata
	readCalled     bool
}

func (m *mockMetadataReader) Read(path string) (controller.FileMetadata, error) {
	m.readCalled = true
	return m.readReturnData, nil
}

type mockFileCacheDataStorer struct {
	uniqueId          string
	metadataCalled    bool
	modtimeReturnData time.Time
}

func (m *mockFileCacheDataStorer) Modtime(path string) (time.Time, error) {
	return m.modtimeReturnData, nil
}

func (m *mockFileCacheDataStorer) Metadata(path string) (controller.FileMetadata, error) {
	m.metadataCalled = true
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
	statReturnValue fs.FileInfo
}

func (m *mockStater) Stat(name string) (fs.FileInfo, error) {
	return m.statReturnValue, nil
}

type fileInfo struct {
	modtime time.Time
}

func (f *fileInfo) Name() string {
	return ""
}

func (f *fileInfo) Size() int64 {
	return 0
}

func (f *fileInfo) Mode() os.FileMode {
	return 0777
}

func (f *fileInfo) ModTime() time.Time {
	return f.modtime
}

func (f *fileInfo) IsDir() bool {
	return false
}

func (f *fileInfo) Sys() interface{} {
	return nil
}
