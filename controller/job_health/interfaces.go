package job_health

import "time"

type NowSincer interface {
	Now() time.Time
	Since(time.Time) time.Duration
}
