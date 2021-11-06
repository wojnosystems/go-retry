package retry

import (
	"context"
	"github.com/wojnosystems/go-retry/core"
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

func (c *Forever) Retry(ctx context.Context, cb core.CallbackFunc) (err error) {
	return core.LoopForever(ctx, cb, func(_ uint64) {
		core.Sleep(ctx, c.WaitBetweenAttempts)
	})
}
