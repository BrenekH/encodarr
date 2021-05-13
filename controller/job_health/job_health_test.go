package job_health

import "testing"

func TestImplement(t *testing.T) {
	t.Errorf("test file not implemented")
}

// TODO: Implement
// Tests to create
//   - Run calls time.Since() and SettingsStorer.HealthCheckInterval()
//   - Run only calls DataStorer.DispatchedJobs() when time.Since is greater than SettingsStorer.HealthCheckInterval()
//   - Various scenarios around dispatched jobs LastUpdated field being higher or lower than HealthCheckTimeout
