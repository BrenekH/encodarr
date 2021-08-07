package jobhealth

import "time"

type nowSincer interface {
	Now() time.Time
	Since(time.Time) time.Duration
}
