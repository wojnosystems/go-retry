package retry

import "time"

// LinearUpTo will retry waiting an exponentially increasing time
// between each request until MaxAttempts is reached.
type LinearUpTo struct {
	InitialWaitBetweenAttempts time.Duration
	GrowthFactor               float64
	MaxAttempts                uint
}

func (c *LinearUpTo) Retry(cb func() (err error)) (err error) {
	return loopUpTo(cb, func(i uint64) {
		sleepTime := linearSleepTime(c.InitialWaitBetweenAttempts, c.GrowthFactor, i)
		time.Sleep(sleepTime)
	}, uint64(c.MaxAttempts))
}
