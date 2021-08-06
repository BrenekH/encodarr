package job_health

import "time"

type TimeNowSince struct{}

func (t TimeNowSince) Now() time.Time {
	return time.Now()
}

func (t TimeNowSince) Since(tt time.Time) time.Duration {
	return time.Since(tt)
}
