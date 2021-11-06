package retryLoop

import (
	"context"
)

func forever(_ uint64) bool {
	return true
}

// Forever will continuously call the callback (cb) until it succeeds or
// returns a non-retryable error
func Forever(ctx context.Context, cb CallbackFunc, wait WaitBetweenAttemptsFunc) (err error) {
	return Until(ctx, cb, wait, forever)
}
