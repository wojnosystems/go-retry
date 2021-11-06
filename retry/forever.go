package retry

import (
	"context"
	"github.com/wojnosystems/go-retry/retryLoop"
	"github.com/wojnosystems/go-retry/retrySleep"
	"time"
)

// Forever will retry forever until your call succeeds or a non-retryable error is reported
type Forever struct {
	WaitBetweenAttempts time.Duration
}

func NewForever(
	waitBetweenAttempts time.Duration,
) *Forever {
	return &Forever{
		WaitBetweenAttempts: waitBetweenAttempts,
	}
}

func (c *Forever) Retry(ctx context.Context, cb retryLoop.CallbackFunc) (err error) {
	return retryLoop.Forever(ctx, cb, func(_ uint64) {
		retrySleep.WithContext(ctx, c.WaitBetweenAttempts)
	})
}
