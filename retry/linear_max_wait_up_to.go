package retry

import (
	"context"
	"github.com/wojnosystems/go-retry/retryLoop"
	"github.com/wojnosystems/go-retry/retrySleep"
	"time"
)

// LinearMaxWaitUpTo retries, but multiplies the InitialWaitBetweenAttempts by the growth factor and adds it to the
// previous wait time each attempt. This allows for an un-bounded linearly-growing backoff.
// BackoffTime(i) = InitialWaitBetweenAttempts * (1 + GrowthFactor*i)
// This is just like Linear and LinearUpTo, but also adds in a cap on the time spent waiting between requests
type LinearMaxWaitUpTo struct {
	InitialWaitBetweenAttempts time.Duration
	GrowthFactor               float64
	MaxAttempts                uint
	MaxWaitBetweenAttempts     time.Duration
}

func NewLinearMaxWaitUpTo(
	initialWaitBetweenAttempts time.Duration,
	growthFactor float64,
	maxAttempts uint,
	maxWaitBetweenAttempts time.Duration,
) *LinearMaxWaitUpTo {
	return &LinearMaxWaitUpTo{
		InitialWaitBetweenAttempts: initialWaitBetweenAttempts,
		GrowthFactor:               growthFactor,
		MaxAttempts:                maxAttempts,
		MaxWaitBetweenAttempts:     maxWaitBetweenAttempts,
	}
}

func (c *LinearMaxWaitUpTo) Retry(ctx context.Context, cb retryLoop.CallbackFunc) (err error) {
	return retryLoop.UpTo(ctx, cb, func(i uint64) {
		sleepTime := linearSleepTime(c.InitialWaitBetweenAttempts, c.GrowthFactor, i)
		retrySleep.WithContext(ctx, minDuration(sleepTime, c.MaxWaitBetweenAttempts))
	}, uint64(c.MaxAttempts))
}
