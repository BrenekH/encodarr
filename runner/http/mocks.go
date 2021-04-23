package http

import (
	"io"
	netHTTP "net/http"
	"time"
)

type mockHTTPClient struct {
	DoResponse  netHTTP.Response
	LastRequest netHTTP.Request

	doCalled bool
}

func (m *mockHTTPClient) Do(req *netHTTP.Request) (*netHTTP.Response, error) {
	m.doCalled = true
	m.LastRequest = *req
	return &m.DoResponse, nil
}

// mockCurrentTime is a mock struct for the CurrentTimer interface.
type mockCurrentTime struct {
	time time.Time
}

func (m *mockCurrentTime) Now() time.Time {
	return m.time
}

// mockFS is a mock struct for the FSer interface.
type mockFS struct {
	createdFiles []string
	openedFiles  []string
}

func (m *mockFS) Create(name string) (Filer, error) {
	if m.createdFiles == nil {
		m.createdFiles = make([]string, 1)
	}
	m.createdFiles = append(m.createdFiles, name)
	return &mockFiler{name}, nil
}

func (m *mockFS) Open(name string) (Filer, error) {
	if m.openedFiles == nil {
		m.openedFiles = make([]string, 1)
	}
	m.openedFiles = append(m.openedFiles, name)
	return &mockFiler{name}, nil
}

// mockFiler is a mock struct for the Filer interface.
type mockFiler struct {
	name string
}

func (m *mockFiler) Close() error { return nil }

func (m *mockFiler) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (m *mockFiler) Write(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (m *mockFiler) Name() string {
	return m.name
}
