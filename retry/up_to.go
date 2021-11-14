package retry

import (
	"context"
	"github.com/wojnosystems/go-retry/retryLoop"
	"github.com/wojnosystems/go-retry/retrySleep"
	"time"
)

// UpTo retries up to MaxAttempts and waits the same WaitBetweenAttempts duration between each retryable error.
type UpTo struct {
	retryStrategy
	// WaitBetweenAttempts
	WaitBetweenAttempts time.Duration

	// MaxAttempts is how many failed tries to attempt before returning an error and giving up
	MaxAttempts uint
}

func NewUpTo(
	waitBetweenAttempts time.Duration,
	maxAttempts uint,
) *UpTo {
	return &UpTo{
		WaitBetweenAttempts: waitBetweenAttempts,
		MaxAttempts:         maxAttempts,
	}
}

func (c *UpTo) Retry(ctx context.Context, cb retryLoop.CallbackFunc) (err error) {
	return retryLoop.UpTo(ctx, cb, func(i uint64) {
		retrySleep.WithContext(ctx, c.WaitBetweenAttempts)
	}, uint64(c.MaxAttempts))
}
