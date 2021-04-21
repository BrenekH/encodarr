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
type mockFS struct{}

func (m *mockFS) Create(name string) (Filer, error) {
	return &mockFiler{}, nil
}

func (m *mockFS) Open(name string) (Filer, error) {
	return &mockFiler{}, nil
}

// mockFiler is a mock struct for the Filer interface.
type mockFiler struct{}

func (m *mockFiler) Close() error { return nil }

func (m *mockFiler) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (m *mockFiler) Write(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (m *mockFiler) Name() string {
	return "mockFile.mkv"
}
