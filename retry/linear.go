package retry

import (
	"context"
	"github.com/wojnosystems/go-retry/retryLoop"
	"github.com/wojnosystems/go-retry/retrySleep"
	"time"
)

// Linear retries, but multiplies the InitialWaitBetweenAttempts by the growth factor and adds it to the
// previous wait time each attempt. This allows for an un-bounded linearly-growing backoff.
// BackoffTime(i) = InitialWaitBetweenAttempts * (1 + GrowthFactor*i)
type Linear struct {
	retryStrategy
	InitialWaitBetweenAttempts time.Duration
	GrowthFactor               float64
}

func NewLinear(initialWaitBetweenAttempts time.Duration, growthFactor float64) *Linear {
	return &Linear{
		InitialWaitBetweenAttempts: initialWaitBetweenAttempts,
		GrowthFactor:               growthFactor,
	}
}

func (c *Linear) Retry(ctx context.Context, cb retryLoop.CallbackFunc) (err error) {
	return retryLoop.Forever(ctx, cb, func(i uint64) {
		sleepTime := linearSleepTime(c.InitialWaitBetweenAttempts, c.GrowthFactor, i)
		retrySleep.WithContext(ctx, sleepTime)
	})
}
