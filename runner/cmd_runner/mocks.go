package cmd_runner

import (
	"bytes"
	"io"
	"time"
)

type mockSincer struct {
	t time.Time
}

func (m *mockSincer) Since(t time.Time) time.Duration {
	return m.t.Sub(t)
}

type mockCommander struct {
	cmder mockCmder

	lastCallArgs []string
}

func (m *mockCommander) Command(name string, args ...string) Cmder {
	m.lastCallArgs = make([]string, 0)
	m.lastCallArgs = append(m.lastCallArgs, args...)
	return &m.cmder
}

type mockCmder struct {
	statusCode int
}

func (m *mockCmder) Start() error {
	return nil
}

func (m *mockCmder) StderrPipe() (io.ReadCloser, error) {
	return io.NopCloser(&bytes.Buffer{}), io.EOF
}

func (m *mockCmder) Wait() error {
	if m.statusCode == 0 {
		return nil
	}
	return mockExitError{statusCode: m.statusCode}
}

type mockExitError struct {
	statusCode int
}

func (m mockExitError) Error() string {
	return "mock exit error"
}

func (m mockExitError) ExitCode() int {
	return m.statusCode
}
