package retry

import (
	"context"
	"github.com/wojnosystems/go-retry/retryLoop"
	"github.com/wojnosystems/go-retry/retrySleep"
	"time"
)

// ExponentialMaxWaitUpTo retries but backs off exponentially by the formula:
// BackoffTime(i) = InitialWaitBetweenAttempts * (1 + GrowthFactor)^i
// where i [0,INF) and represents the number of times we've delayed after a failed attempt before
// This is just like Exponential and ExponentialUpTo, but also adds in a cap on the time spent waiting between requests
type ExponentialMaxWaitUpTo struct {
	retryStrategy
	InitialWaitBetweenAttempts time.Duration
	GrowthFactor               float64
	MaxAttempts                uint
	MaxWaitBetweenAttempts     time.Duration
}

func NewExponentialMaxWaitUpTo(
	initialWaitBetweenAttempts time.Duration,
	growthFactor float64,
	maxAttempts uint,
	maxWaitBetweenAttempts time.Duration,
) *ExponentialMaxWaitUpTo {
	return &ExponentialMaxWaitUpTo{
		InitialWaitBetweenAttempts: initialWaitBetweenAttempts,
		GrowthFactor:               growthFactor,
		MaxAttempts:                maxAttempts,
		MaxWaitBetweenAttempts:     maxWaitBetweenAttempts,
	}
}

func (c *ExponentialMaxWaitUpTo) Retry(ctx context.Context, cb retryLoop.CallbackFunc) (err error) {
	return retryLoop.UpTo(ctx, cb, func(i uint64) {
		sleepTime := exponentialSleepTime(c.InitialWaitBetweenAttempts, c.GrowthFactor, i)
		retrySleep.WithContext(ctx, minDuration(sleepTime, c.MaxWaitBetweenAttempts))
	}, uint64(c.MaxAttempts))
}
