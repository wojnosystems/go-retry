package retry

import (
	"context"
	"github.com/wojnosystems/go-retry/core"
	"time"
)

// UpTo retries up to MaxAttempts and waits the same WaitBetweenAttempts duration between each retryable error.
type UpTo struct {
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

func (c *UpTo) Retry(ctx context.Context, cb core.CallbackFunc) (err error) {
	return core.LoopUpTo(ctx, cb, func(i uint64) {
		core.Sleep(ctx, c.WaitBetweenAttempts)
	}, uint64(c.MaxAttempts))
}
