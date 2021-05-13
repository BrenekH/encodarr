package settings

import (
	"bytes"
	"io"
)

type mockReadWriteCloser struct {
	readCalled  bool
	writeCalled bool
	closeCalled bool

	bR *bytes.Reader
}

func (m *mockReadWriteCloser) Read(p []byte) (n int, err error) {
	if !m.readCalled {
		m.bR = bytes.NewReader([]byte("{}"))
	}
	m.readCalled = true
	return m.bR.Read(p)
}

func (m *mockReadWriteCloser) Write(p []byte) (n int, err error) {
	m.writeCalled = true
	return io.Discard.Write(p)
}

func (m *mockReadWriteCloser) Close() error {
	m.closeCalled = true
	return nil
}
