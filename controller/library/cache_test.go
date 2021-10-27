package library

import (
	"errors"
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
		t.Error("FileCacheDataStorer.Metadata() was not called even though it should have been")
	}

	if m.readCalled {
		t.Error("MetadataReader.Read() was called even though it should not have been")
	}
}

func TestCacheReadStaterErr(t *testing.T) {
	mtime := time.Unix(10000, 100)

	m := mockMetadataReader{}
	f := mockFileCacheDataStorer{modtimeReturnData: mtime}
	l := mockLogger{}

	cache := NewCache(&m, &f, &l)
	s := mockStater{statReturnValue: &fileInfo{modtime: mtime}, err: errors.New("some error")}
	cache.stater = &s

	cache.Read("test")

	if !m.readCalled {
		t.Error("MetadataReader.Read() was not called even though it should have been")
	}
}

func TestCacheReadModtimeErr(t *testing.T) {
	mtime := time.Unix(10000, 100)

	m := mockMetadataReader{}
	f := mockFileCacheDataStorer{
		modtimeReturnData: mtime,
		modtimeErr:        errors.New("some error"),
	}
	l := mockLogger{}

	cache := NewCache(&m, &f, &l)
	s := mockStater{statReturnValue: &fileInfo{modtime: mtime}}
	cache.stater = &s

	cache.Read("test")

	if !f.modtimeCalled {
		t.Error("FileCacheDataStorer.Modtime() was not called even though it should have been")
	}

	if !m.readCalled {
		t.Error("MetadataReader.Read() was not called even though it should have been")
	}
}

func TestCacheReadMetadataErr(t *testing.T) {
	mtime := time.Unix(10000, 100)

	m := mockMetadataReader{}
	f := mockFileCacheDataStorer{
		modtimeReturnData: mtime,
		metadataErr:       errors.New("some error"),
	}
	l := mockLogger{}

	cache := NewCache(&m, &f, &l)
	s := mockStater{statReturnValue: &fileInfo{modtime: mtime}}
	cache.stater = &s

	cache.Read("test")

	if !f.metadataCalled {
		t.Error("FileCacheDataStorer.Metadata() was not called even though it should have been")
	}

	if !m.readCalled {
		t.Error("MetadataReader.Read() was not called even though it should have been")
	}
}

func TestCacheReadSaveMetadataErr(t *testing.T) {
	mtime := time.Unix(10000, 100)

	m := mockMetadataReader{}
	f := mockFileCacheDataStorer{
		modtimeReturnData: mtime,
		saveMetadataErr:   errors.New("some error"),
	}
	l := mockLogger{}

	cache := NewCache(&m, &f, &l)
	s := mockStater{
		statReturnValue: &fileInfo{modtime: mtime.Add(time.Duration(-10) * time.Minute)},
	}
	cache.stater = &s

	cache.Read("test")

	if !m.readCalled {
		t.Error("MetadataReader.Read() was not called even though it should have been")
	}

	if !f.saveMetadataCalled {
		t.Error("FileCacheDataStorer.SaveMetadata() was not called even though it should have been")
	}
}

func TestCacheReadSaveModTimeErr(t *testing.T) {
	mtime := time.Unix(10000, 100)

	m := mockMetadataReader{}
	f := mockFileCacheDataStorer{
		modtimeReturnData: mtime,
		saveModTimeErr:    errors.New("some error"),
	}
	l := mockLogger{}

	cache := NewCache(&m, &f, &l)
	s := mockStater{
		statReturnValue: &fileInfo{modtime: mtime.Add(time.Duration(-10) * time.Minute)},
	}
	cache.stater = &s

	cache.Read("test")

	if !m.readCalled {
		t.Error("MetadataReader.Read() was not called even though it should have been")
	}

	if !f.saveModTimeCalled {
		t.Error("FileCacheDataStorer.SaveModTime() was not called even though it should have been")
	}
}

func TestStat(t *testing.T) {
	m := mockMetadataReader{}
	f := mockFileCacheDataStorer{}
	l := mockLogger{}

	c := NewCache(&m, &f, &l)

	_, got := c.stater.Stat("some file")

	if got == nil {
		t.Errorf("got %v, but want %v", got, "not nil")
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
	uniqueId           string
	metadataCalled     bool
	modtimeCalled      bool
	saveModTimeCalled  bool
	saveMetadataCalled bool
	metadataErr        error
	modtimeErr         error
	saveModTimeErr     error
	saveMetadataErr    error
	modtimeReturnData  time.Time
}

func (m *mockFileCacheDataStorer) Modtime(path string) (time.Time, error) {
	m.modtimeCalled = true
	return m.modtimeReturnData, m.modtimeErr
}

func (m *mockFileCacheDataStorer) Metadata(path string) (controller.FileMetadata, error) {
	m.metadataCalled = true
	return controller.FileMetadata{}, m.metadataErr
}

func (m *mockFileCacheDataStorer) SaveModtime(path string, t time.Time) error {
	m.saveModTimeCalled = true
	return m.saveModTimeErr
}

func (m *mockFileCacheDataStorer) SaveMetadata(path string, f controller.FileMetadata) error {
	m.saveMetadataCalled = true
	return m.saveMetadataErr
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
	err             error
}

func (m *mockStater) Stat(name string) (fs.FileInfo, error) {
	return m.statReturnValue, m.err
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
