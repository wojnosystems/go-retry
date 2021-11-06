package retry

import (
	"context"
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

func NewExponentialUpTo(
	initialWaitBetweenAttempts time.Duration,
	growthFactor float64,
	maxAttempts uint,
) *ExponentialUpTo {
	return &ExponentialUpTo{
		InitialWaitBetweenAttempts: initialWaitBetweenAttempts,
		GrowthFactor:               growthFactor,
		MaxAttempts:                maxAttempts,
	}
}

func (c *ExponentialUpTo) Retry(ctx context.Context, cb core.CallbackFunc) (err error) {
	return core.LoopUpTo(ctx, cb, func(i uint64) {
		sleepTime := exponentialSleepTime(c.InitialWaitBetweenAttempts, c.GrowthFactor, i)
		core.Sleep(ctx, sleepTime)
	}, uint64(c.MaxAttempts))
}
