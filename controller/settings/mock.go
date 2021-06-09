package settings

import (
	"bytes"
	"io"
)

type mockReadWriteSeekCloser struct {
	readCalled  bool
	writeCalled bool
	seekCalled  bool
	closeCalled bool

	bR *bytes.Reader
}

func (m *mockReadWriteSeekCloser) Read(p []byte) (int, error) {
	if !m.readCalled {
		m.bR = bytes.NewReader([]byte("{}"))
	}
	m.readCalled = true
	return m.bR.Read(p)
}

func (m *mockReadWriteSeekCloser) Write(p []byte) (int, error) {
	m.writeCalled = true
	return io.Discard.Write(p)
}

func (m *mockReadWriteSeekCloser) Seek(offset int64, whence int) (int64, error) {
	m.seekCalled = true
	if !m.readCalled {
		return 0, nil
	}
	return m.bR.Seek(offset, whence)
}

func (m *mockReadWriteSeekCloser) Truncate(size int64) error {
	m.readCalled = false // Effectively resets the byteReader
	return nil
}

func (m *mockReadWriteSeekCloser) Close() error {
	m.closeCalled = true
	return nil
}
