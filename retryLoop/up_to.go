package retryLoop

import (
	"context"
)

// UpTo will call callback until it returns a non-retryable error, success, or maxAttempts is exceeded
func UpTo(ctx context.Context, callback CallbackFunc, wait WaitBetweenAttemptsFunc, maxAttempts uint64) (err error) {
	return Until(ctx, callback, wait, func(timesAttempted uint64) bool {
		return timesAttempted < maxAttempts
	})
}
