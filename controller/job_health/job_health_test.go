package job_health

import "testing"

// Run calls time.Since() and SettingsStorer.HealthCheckInterval()
func TestTimeSinceAndSSHealthCheckIntervalCalled(t *testing.T) {
	ds := mockDataStorer{}
	ss := mockSettingsStorer{}
	c := NewChecker(&ds, &ss)

	mNS := mockNowSincer{}
	c.nowSincer = &mNS

	c.Run()

	if !mNS.sinceCalled {
		t.Errorf("expected NowSincer.Since() to be called")
	}
	if !ss.healthCheckIntCalled {
		t.Errorf("expected SettingsStorer.HealthCheckInterval() to be called")
	}
}

// TODO: Implement
// Tests to create
//   - Run only calls DataStorer.DispatchedJobs() when time.Since is greater than SettingsStorer.HealthCheckInterval()
//   - Various scenarios around dispatched jobs LastUpdated field being higher or lower than HealthCheckTimeout
