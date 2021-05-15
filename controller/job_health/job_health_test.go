package job_health

import (
	"testing"
	"time"

	"github.com/BrenekH/encodarr/controller"
)

// Run calls time.Since() and SettingsStorer.HealthCheckInterval()
func TestTimeSinceAndSSHealthCheckIntervalCalled(t *testing.T) {
	ds := mockDataStorer{}
	ss := mockSettingsStorer{}
	c := NewChecker(&ds, &ss, &mockLogger{})

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
			c := NewChecker(&ds, &ss, &mockLogger{})

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

// Various scenarios around dispatched jobs LastUpdated field being higher or lower than HealthCheckTimeout
func TestCorrectNullingBehavior(t *testing.T) {
	tests := []struct {
		name                    string
		healthCheckTimeout      uint64
		jobLastUpdatedDur       time.Duration // Sets what time.Since returns as the duration on the second call
		expectUUIDToBeNullified bool
	}{
		{
			name:                    "DJob.LastUpdated is smaller than the timeout",
			healthCheckTimeout:      uint64(time.Second * 32),
			jobLastUpdatedDur:       time.Second * 16,
			expectUUIDToBeNullified: false,
		},
		{
			name:                    "DJob.LastUpdated is equal to the timeout",
			healthCheckTimeout:      uint64(time.Second * 32),
			jobLastUpdatedDur:       time.Second * 32,
			expectUUIDToBeNullified: true,
		},
		{
			name:                    "DJob.LastUpdated is larger than the timeout",
			healthCheckTimeout:      uint64(time.Second * 32),
			jobLastUpdatedDur:       time.Second * 48,
			expectUUIDToBeNullified: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ds := mockDataStorer{
				dJobs: []controller.DispatchedJob{
					{
						UUID:        "test",
						Runner:      "TestRunner",
						Job:         controller.Job{},
						Status:      controller.JobStatus{},
						LastUpdated: time.Unix(0, 0),
					},
				},
			}
			ss := mockSettingsStorer{
				healthCheckInt:     uint64(time.Second * 1),
				healthCheckTimeout: test.healthCheckTimeout,
			}
			c := NewChecker(&ds, &ss, &mockLogger{})

			mNS := mockNowSincer{
				sinceResp:  time.Second * 2,
				sinceResp2: test.jobLastUpdatedDur,
			}
			c.nowSincer = &mNS

			nulledUUIDs := c.Run()

			if !ds.dJobsCalled {
				t.Errorf("expected DataStorer.DispatchedJobs() to be called")
				return
			}

			if mNS.sinceTimesCalled != 2 {
				t.Errorf("expected NowSincer.Since() to be called twice but it was called %v times", mNS.sinceTimesCalled)
				return
			}

			if len(nulledUUIDs) > 0 && !test.expectUUIDToBeNullified {
				t.Errorf("received a nullified UUID when one wasn't expected")
			}
		})
	}
}

// TODO: Test DataStorer.DeleteJob error handling
