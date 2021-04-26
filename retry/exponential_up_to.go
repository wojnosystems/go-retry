package retry

import "time"

// ExponentialUpTo will retry waiting an exponentially increasing time
// between each request until MaxAttempts is reached.
type ExponentialUpTo struct {
	InitialWaitBetweenAttempts time.Duration
	GrowthFactor               float64
	MaxAttempts                uint
}

func (c *ExponentialUpTo) Retry(cb func() (err error)) (err error) {
	return loopUpTo(cb, func(i uint64) {
		sleepTime := exponentialSleepTime(c.InitialWaitBetweenAttempts, c.GrowthFactor, i)
		time.Sleep(sleepTime)
	}, uint64(c.MaxAttempts))
}
