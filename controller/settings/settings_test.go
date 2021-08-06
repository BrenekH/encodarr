package settings

import (
	"testing"

	"github.com/BrenekH/encodarr/controller"
)

func TestFileReadCalled(t *testing.T) {
	ss := defaultSettings()
	mRWC := mockReadWriteSeekCloser{}
	ss.file = &mRWC

	if err := ss.Load(); err != nil {
		t.Errorf("unexpected error from SettingsStore.Read(): %v", err)
	}

	if !mRWC.readCalled {
		t.Errorf("expected SettingsStore.file.Read() to be called")
	}
}

func TestFileWriteCalled(t *testing.T) {
	ss := defaultSettings()
	mRWC := mockReadWriteSeekCloser{}
	ss.file = &mRWC

	if err := ss.Save(); err != nil {
		t.Errorf("unexpected error from SettingsStore.Save(): %v", err)
	}

	if !mRWC.writeCalled {
		t.Errorf("expected SettingsStore.file.Write() to be called")
	}
}

func TestFileCloseCalled(t *testing.T) {
	ss := defaultSettings()
	mRWC := mockReadWriteSeekCloser{}
	ss.file = &mRWC

	if err := ss.Close(); err != nil {
		t.Errorf("unexpected error from SettingsStore.Close(): %v", err)
	}

	if !mRWC.closeCalled {
		t.Errorf("expected SettingsStore.file.Close() to be called")
	}
}

func TestCloseSetsClosed(t *testing.T) {
	ss := defaultSettings()
	mRWC := mockReadWriteSeekCloser{}
	ss.file = &mRWC

	if err := ss.Close(); err != nil {
		t.Errorf("unexpected error from SettingsStore.Close(): %v", err)
	}

	if !ss.closed {
		t.Errorf("expected SettingsStore.closed to be true, but it was false")
	}
}

func TestErrorReturnedAfterFileIsClosed(t *testing.T) {
	ss := defaultSettings()
	mRWC := mockReadWriteSeekCloser{}
	ss.file = &mRWC

	if err := ss.Close(); err != nil {
		t.Errorf("unexpected error from SettingsStore.Close(): %v", err)
	}

	if err := ss.Load(); err != controller.ErrClosed {
		t.Errorf("expected controller.ErrClosed error from Load() after Close() is called, but got %v instead", err)
	}

	if err := ss.Save(); err != controller.ErrClosed {
		t.Errorf("expected controller.ErrClosed error from Save() after Close() is called, but got %v instead", err)
	}
}

// TODO: Test that file is truncated on call to ss.Save
