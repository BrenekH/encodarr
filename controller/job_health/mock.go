package job_health

import (
	"time"

	"github.com/BrenekH/encodarr/controller"
)

type mockNowSincer struct {
	nowResp   time.Time
	sinceResp time.Duration

	nowCalled   bool
	sinceCalled bool
}

func (m *mockNowSincer) Now() time.Time {
	m.nowCalled = true
	return m.nowResp
}

func (m *mockNowSincer) Since(time.Time) time.Duration {
	m.sinceCalled = true
	return m.sinceResp
}

type mockDataStorer struct {
	dJobsCalled bool
}

func (m *mockDataStorer) DispatchedJobs() (djs []controller.DispatchedJob) {
	m.dJobsCalled = true
	return
}
func (m *mockDataStorer) DeleteJob(uuid controller.UUID) {}

type mockSettingsStorer struct {
	healthCheckIntCalled bool

	healthCheckInt uint64
}

func (m *mockSettingsStorer) HealthCheckInterval() uint64 {
	m.healthCheckIntCalled = true
	return m.healthCheckInt
}

func (m *mockSettingsStorer) Load() (err error)              { return }
func (m *mockSettingsStorer) Save() (err error)              { return }
func (m *mockSettingsStorer) Close() (err error)             { return }
func (m *mockSettingsStorer) SetHealthCheckInterval(uint64)  {}
func (m *mockSettingsStorer) HealthCheckTimeout() (u uint64) { return }
func (m *mockSettingsStorer) SetHealthCheckTimeout(uint64)   {}
func (m *mockSettingsStorer) LogVerbosity() (s string)       { return }
func (m *mockSettingsStorer) SetLogVerbosity(string)         {}
