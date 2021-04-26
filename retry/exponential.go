package retry

import (
	"math"
	"time"
)

// Exponential retries but backs off exponentially by the formula:
// BackoffTime(i) = InitialWaitBetweenAttempts * (1 + GrowthFactor)^i
// where i [0,)
type Exponential struct {
	InitialWaitBetweenAttempts time.Duration
	GrowthFactor               float64
}

func (c *Exponential) Retry(cb func() (err error)) (err error) {
	i := uint64(0)
	return loopForever(cb, func() {
		sleepTime := exponentialSleepTime(c.InitialWaitBetweenAttempts, c.GrowthFactor, i)
		if i < math.MaxUint64 {
			i++
		}
		time.Sleep(sleepTime)
	})
}

func exponentialSleepTime(initial time.Duration, growthFactor float64, iteration uint64) time.Duration {
	return time.Duration(
		float64(initial) * math.Pow(1.0+growthFactor, float64(iteration)))
}
