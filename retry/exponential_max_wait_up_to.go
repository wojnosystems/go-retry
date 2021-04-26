package retry

import "time"

// ExponentialMaxWaitUpTo will retry waiting an exponentially increasing time
// between each request until MaxAttempts is reached. The wait time will stop getting larger than MaxWaitBetweenAttempts
type ExponentialMaxWaitUpTo struct {
	InitialWaitBetweenAttempts time.Duration
	GrowthFactor               float64
	MaxAttempts                uint
	MaxWaitBetweenAttempts     time.Duration
}

func (c *ExponentialMaxWaitUpTo) Retry(cb func() (err error)) (err error) {
	return loopUpTo(cb, func(i uint64) {
		sleepTime := exponentialSleepTime(c.InitialWaitBetweenAttempts, c.GrowthFactor, i)
		time.Sleep(minDuration(sleepTime, c.MaxWaitBetweenAttempts))
	}, uint64(c.MaxAttempts))
}
