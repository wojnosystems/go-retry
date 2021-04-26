package retry

import (
	"math"
	"time"
)

// Linear retries, but multiplies the InitialWaitBetweenAttempts by the growth factor and adds it to the
// previous wait time each attempt. This allows for an un-bounded linearly-growing backoff.
type Linear struct {
	InitialWaitBetweenAttempts time.Duration
	GrowthFactor               float64
}

func (c *Linear) Retry(cb func() (err error)) (err error) {
	i := uint64(0)
	return loopForever(cb, func() {
		sleepTime := linearSleepTime(c.InitialWaitBetweenAttempts, c.GrowthFactor, i)
		if i < math.MaxUint64 {
			i++
		}
		time.Sleep(sleepTime)
	})
}

func linearSleepTime(initial time.Duration, growthFactor float64, iteration uint64) time.Duration {
	return time.Duration(float64(initial) + (float64(initial) * growthFactor * float64(iteration)))
}
