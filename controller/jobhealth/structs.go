package jobhealth

import "time"

type timeNowSince struct{}

func (t timeNowSince) Now() time.Time {
	return time.Now()
}

func (t timeNowSince) Since(tt time.Time) time.Duration {
	return time.Since(tt)
}
