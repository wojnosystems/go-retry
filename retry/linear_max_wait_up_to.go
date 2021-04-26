package retry

import "time"

// LinearMaxWaitUpTo will retry waiting an exponentially increasing time
// between each request until MaxAttempts is reached. However, it will stop growing the duration between requests after MaxWaitBetweenAttempts is reached
type LinearMaxWaitUpTo struct {
	InitialWaitBetweenAttempts time.Duration
	GrowthFactor               float64
	MaxAttempts                uint
	MaxWaitBetweenAttempts     time.Duration
}

func (c *LinearMaxWaitUpTo) Retry(cb func() (err error)) (err error) {
	return loopUpTo(cb, func(i uint64) {
		sleepTime := linearSleepTime(c.InitialWaitBetweenAttempts, c.GrowthFactor, i)
		time.Sleep(minDuration(sleepTime, c.MaxWaitBetweenAttempts))
	}, uint64(c.MaxAttempts))
}

func minDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
