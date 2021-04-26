package retry

import (
	"github.com/wojnosystems/go-retry/core"
	"time"
)

// Linear retries, but multiplies the InitialWaitBetweenAttempts by the growth factor and adds it to the
// previous wait time each attempt. This allows for an un-bounded linearly-growing backoff.
// BackoffTime(i) = InitialWaitBetweenAttempts * (1 + GrowthFactor*i)
type Linear struct {
	InitialWaitBetweenAttempts time.Duration
	GrowthFactor               float64
}

func (c *Linear) Retry(cb core.CallbackFunc) (err error) {
	return core.LoopForever(cb, func(i uint64) {
		sleepTime := linearSleepTime(c.InitialWaitBetweenAttempts, c.GrowthFactor, i)
		time.Sleep(sleepTime)
	})
}

func linearSleepTime(initial time.Duration, growthFactor float64, iteration uint64) time.Duration {
	return time.Duration(float64(initial) + (float64(initial) * growthFactor * float64(iteration)))
}
