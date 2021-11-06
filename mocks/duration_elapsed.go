package mocks

import "time"

func DurationElapsed(timeThis func()) time.Duration {
	startAt := time.Now()
	timeThis()
	endAt := time.Now()
	return endAt.Sub(startAt)
}
