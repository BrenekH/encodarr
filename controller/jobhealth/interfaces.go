package jobhealth

import "time"

type NowSincer interface {
	Now() time.Time
	Since(time.Time) time.Duration
}
