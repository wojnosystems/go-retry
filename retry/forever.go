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

func (c *Forever) Retry(ctx context.Context, cb core.CallbackFunc) (err error) {
	return core.LoopForever(cb, func(_ uint64) {
		core.Sleep(ctx, c.WaitBetweenAttempts)
	})
}
