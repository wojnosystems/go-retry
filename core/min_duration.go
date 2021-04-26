package core

import "time"

// MinDuration returns the minimum of a and b
func MinDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
