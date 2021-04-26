package retry

import (
	"github.com/wojnosystems/go-retry/core"
	"time"
)

// ExponentialUpTo retries but backs off exponentially by the formula:
// BackoffTime(i) = InitialWaitBetweenAttempts * (1 + GrowthFactor)^i
// where i [0,INF) and represents the number of times we've delayed after a failed attempt before
// This is just like Exponential, except that it will also only execute a finite number of times before stopping
type ExponentialUpTo struct {
	InitialWaitBetweenAttempts time.Duration
	GrowthFactor               float64
	MaxAttempts                uint
}

func (c *ExponentialUpTo) Retry(cb core.CallbackFunc) (err error) {
	return core.LoopUpTo(cb, func(i uint64) {
		sleepTime := exponentialSleepTime(c.InitialWaitBetweenAttempts, c.GrowthFactor, i)
		time.Sleep(sleepTime)
	}, uint64(c.MaxAttempts))
}
