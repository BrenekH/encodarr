package jobhealth

import (
	"fmt"
	"time"

	"github.com/BrenekH/encodarr/controller"
)

type mockNowSincer struct {
	nowResp    time.Time
	sinceResp  time.Duration
	sinceResp2 time.Duration

	nowCalled        bool
	sinceCalled      bool
	sinceTimesCalled int
}

func (m *mockNowSincer) Now() time.Time {
	m.nowCalled = true
	return m.nowResp
}

func (m *mockNowSincer) Since(time.Time) time.Duration {
	m.sinceTimesCalled++
	if m.sinceCalled {
		return m.sinceResp2
	}
	m.sinceCalled = true
	return m.sinceResp
}

type mockDataStorer struct {
	dJobsCalled bool

	dJobs []controller.DispatchedJob

	deleteErrAmount int
}

func (m *mockDataStorer) DispatchedJobs() []controller.DispatchedJob {
	m.dJobsCalled = true
	return m.dJobs
}

func (m *mockDataStorer) DeleteJob(uuid controller.UUID) error {
	if m.deleteErrAmount == 0 {
		return nil
	}
	m.deleteErrAmount--
	return fmt.Errorf("random error")
}

type mockSettingsStorer struct {
	healthCheckIntCalled bool

	healthCheckInt     uint64
	healthCheckTimeout uint64
}

func (m *mockSettingsStorer) HealthCheckInterval() uint64 {
	m.healthCheckIntCalled = true
	return m.healthCheckInt
}

func (m *mockSettingsStorer) HealthCheckTimeout() uint64 {
	return m.healthCheckTimeout
}

func (m *mockSettingsStorer) Load() (err error)             { return }
func (m *mockSettingsStorer) Save() (err error)             { return }
func (m *mockSettingsStorer) Close() (err error)            { return }
func (m *mockSettingsStorer) SetHealthCheckInterval(uint64) {}
func (m *mockSettingsStorer) SetHealthCheckTimeout(uint64)  {}
func (m *mockSettingsStorer) LogVerbosity() (s string)      { return }
func (m *mockSettingsStorer) SetLogVerbosity(string)        {}

type mockLogger struct{}

func (m *mockLogger) Trace(s string, i ...interface{})    {}
func (m *mockLogger) Debug(s string, i ...interface{})    {}
func (m *mockLogger) Info(s string, i ...interface{})     {}
func (m *mockLogger) Warn(s string, i ...interface{})     {}
func (m *mockLogger) Error(s string, i ...interface{})    {}
func (m *mockLogger) Critical(s string, i ...interface{}) {}
