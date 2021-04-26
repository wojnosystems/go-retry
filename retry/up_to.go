package retry

import (
	"github.com/wojnosystems/go-retry/core"
	"time"
)

// UpTo retries up to MaxAttempts and WaitBetweenAttempts duration between each retryable error
type UpTo struct {
	// WaitBetweenAttempts
	WaitBetweenAttempts time.Duration

	// MaxAttempts is how many failed tries to attempt before returning an error and giving up
	MaxAttempts uint
}

func (c *UpTo) Retry(cb core.CallbackFunc) (err error) {
	return core.LoopUpTo(cb, func(i uint64) {
		time.Sleep(c.WaitBetweenAttempts)
	}, uint64(c.MaxAttempts))
}
