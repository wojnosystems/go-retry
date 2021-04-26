package retry

import (
	"github.com/wojnosystems/go-retry/core"
	"time"
)

// ExponentialMaxWaitUpTo retries but backs off exponentially by the formula:
// BackoffTime(i) = InitialWaitBetweenAttempts * (1 + GrowthFactor)^i
// where i [0,INF) and represents the number of times we've delayed after a failed attempt before
// This is just like Exponential and ExponentialUpTo, but also adds in a cap on the time spent waiting between requests
type ExponentialMaxWaitUpTo struct {
	InitialWaitBetweenAttempts time.Duration
	GrowthFactor               float64
	MaxAttempts                uint
	MaxWaitBetweenAttempts     time.Duration
}

func (c *ExponentialMaxWaitUpTo) Retry(cb core.CallbackFunc) (err error) {
	return core.LoopUpTo(cb, func(i uint64) {
		sleepTime := exponentialSleepTime(c.InitialWaitBetweenAttempts, c.GrowthFactor, i)
		time.Sleep(core.MinDuration(sleepTime, c.MaxWaitBetweenAttempts))
	}, uint64(c.MaxAttempts))
}
