package core

import (
	"context"
)

func loopForever(_ uint64) bool {
	return true
}

// LoopForever will continuously call the callback (cb) until it succeeds or
// returns a non-retryable error
func LoopForever(ctx context.Context, cb CallbackFunc, wait DelayBetweenAttemptsFunc) (err error) {
	return LoopUntil(ctx, cb, wait, loopForever)
}
