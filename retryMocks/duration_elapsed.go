package retryMocks

import "time"

// DurationElapsed returns the time interval from when the timeThis callback started to when it returned
// Useful for testing, but should not be used to actually time anything to accurate precision
func DurationElapsed(timeThis func()) time.Duration {
	startAt := time.Now()
	timeThis()
	endAt := time.Now()
	return endAt.Sub(startAt)
}
