package retry

import (
	"context"
	"github.com/wojnosystems/go-retry/core"
	"time"
)

// LinearUpTo retries, but multiplies the InitialWaitBetweenAttempts by the growth factor and adds it to the
// previous wait time each attempt. This allows for an un-bounded linearly-growing backoff.
// BackoffTime(i) = InitialWaitBetweenAttempts * (1 + GrowthFactor*i)
// This is just like Linear, except that it will also only execute a finite number of times before stopping
type LinearUpTo struct {
	InitialWaitBetweenAttempts time.Duration
	GrowthFactor               float64
	MaxAttempts                uint
}

func NewLinearUpTo(
	initialWaitBetweenAttempts time.Duration,
	growthFactor float64,
	maxAttempts uint,
) *LinearUpTo {
	return &LinearUpTo{
		InitialWaitBetweenAttempts: initialWaitBetweenAttempts,
		GrowthFactor:               growthFactor,
		MaxAttempts:                maxAttempts,
	}
}

func (c *LinearUpTo) Retry(ctx context.Context, cb core.CallbackFunc) (err error) {
	return core.LoopUpTo(ctx, cb, func(i uint64) {
		sleepTime := linearSleepTime(c.InitialWaitBetweenAttempts, c.GrowthFactor, i)
		core.Sleep(ctx, sleepTime)
	}, uint64(c.MaxAttempts))
}
