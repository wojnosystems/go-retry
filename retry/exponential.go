package retry

import (
	"github.com/wojnosystems/go-retry/core"
	"math"
	"time"
)

// Exponential retries but backs off exponentially by the formula:
// BackoffTime(i) = InitialWaitBetweenAttempts * (1 + GrowthFactor)^i
// where i [0,INF) and represents the number of times we've delayed after a failed attempt before
type Exponential struct {
	InitialWaitBetweenAttempts time.Duration
	GrowthFactor               float64
}

func (c *Exponential) Retry(cb core.CallbackFunc) (err error) {
	return core.LoopForever(cb, func(i uint64) {
		sleepTime := exponentialSleepTime(c.InitialWaitBetweenAttempts, c.GrowthFactor, i)
		time.Sleep(sleepTime)
	})
}

func exponentialSleepTime(initial time.Duration, growthFactor float64, iteration uint64) time.Duration {
	return time.Duration(
		float64(initial) * math.Pow(1.0+growthFactor, float64(iteration)))
}
