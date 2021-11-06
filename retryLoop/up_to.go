package retryLoop

import (
	"context"
)

// UpTo will call cb until it returns a non-retryable error or success or maxAttempts is exceeded
func UpTo(ctx context.Context, cb CallbackFunc, wait WaitBetweenAttemptsFunc, maxAttempts uint64) (err error) {
	return Until(ctx, cb, wait, func(timesAttempted uint64) bool {
		return timesAttempted < maxAttempts
	})
}
