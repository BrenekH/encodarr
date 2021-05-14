package job_health

import (
	"testing"
	"time"
)

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

// Run only calls DataStorer.DispatchedJobs() when time.Since is greater than SettingsStorer.HealthCheckInterval()
func TestHealthCheckRunsUnderCorrectConditions(t *testing.T) {
	tests := []struct {
		name             string
		healthCheckInt   uint64
		sinceResp        time.Duration
		djCalledExpected bool
	}{
		{
			name:             "Since returns duration smaller than interval",
			healthCheckInt:   uint64(time.Second * 32),
			sinceResp:        time.Second * 16,
			djCalledExpected: false,
		},
		{
			name:             "Since returns duration equal to interval",
			healthCheckInt:   uint64(time.Second * 32),
			sinceResp:        time.Second * 32,
			djCalledExpected: true,
		},
		{
			name:             "Since returns duration larger than interval",
			healthCheckInt:   uint64(time.Second * 32),
			sinceResp:        time.Second * 48,
			djCalledExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ds := mockDataStorer{}
			ss := mockSettingsStorer{
				healthCheckInt: test.healthCheckInt,
			}
			c := NewChecker(&ds, &ss)

			mNS := mockNowSincer{
				sinceResp: test.sinceResp,
			}
			c.nowSincer = &mNS

			c.Run()

			if ds.dJobsCalled != test.djCalledExpected {
				t.Errorf("expected ds.dJobsCalled to be %v, but it was %v instead", test.djCalledExpected, ds.dJobsCalled)
			}
		})
	}
}

// TODO: Implement
// Tests to create
//   - Various scenarios around dispatched jobs LastUpdated field being higher or lower than HealthCheckTimeout
