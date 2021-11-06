package retry

import "time"

// minDuration returns the minimum of a and b
func minDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
