package cmd_runner

import "time"

type mockSincer struct {
	t time.Time
}

func (m *mockSincer) Since(t time.Time) time.Duration {
	return m.t.Sub(t)
}
