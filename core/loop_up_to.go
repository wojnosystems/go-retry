package core

import (
	"context"
)

// LoopUpTo will call cb until it returns a non-retryable error or success or maxAttempts is exceeded
func LoopUpTo(ctx context.Context, cb CallbackFunc, wait DelayBetweenAttemptsFunc, maxAttempts uint64) (err error) {
	return LoopUntil(ctx, cb, wait, func(timesAttempted uint64) bool {
		return timesAttempted < maxAttempts
	})
}
