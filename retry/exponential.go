package retry

import (
	"context"
	"github.com/wojnosystems/go-retry/retryLoop"
	"github.com/wojnosystems/go-retry/retrySleep"
	"time"
)

// Exponential retries but backs off exponentially by the formula:
// BackoffTime(i) = InitialWaitBetweenAttempts * (1 + GrowthFactor)^i
// where i [0,INF) and represents the number of times we've delayed after a failed attempt before
// should the callback always indicate a retry, this will retry forever
type Exponential struct {
	InitialWaitBetweenAttempts time.Duration
	GrowthFactor               float64
}

func NewExponential(
	initialWaitBetweenAttempts time.Duration,
	growthFactor float64,
) *Exponential {
	return &Exponential{
		InitialWaitBetweenAttempts: initialWaitBetweenAttempts,
		GrowthFactor:               growthFactor,
	}
}

func (c *Exponential) Retry(ctx context.Context, cb retryLoop.CallbackFunc) (err error) {
	return retryLoop.Forever(ctx, cb, func(i uint64) {
		sleepTime := exponentialSleepTime(c.InitialWaitBetweenAttempts, c.GrowthFactor, i)
		retrySleep.WithContext(ctx, sleepTime)
	})
}
