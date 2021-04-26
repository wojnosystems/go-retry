package core

import (
	"github.com/wojnosystems/go-retry/retryAgain"
	"github.com/wojnosystems/go-retry/retryStop"
)

// LoopUpTo will call cb until it returns a non-retryable error or success or maxAttempts is exceeded
func LoopUpTo(cb CallbackFunc, wait DelayBetweenAttemptsFunc, maxAttempts uint64) (err error) {
	i := uint64(0)
	for {
		err = cb()
		if err == retryStop.Success {
			return
		}
		if v, ok := err.(retryAgain.Wrapper); !ok {
			return err
		} else {
			i++
			if i < maxAttempts {
				wait(i - 1)
			} else {
				return v.Err()
			}
		}
	}
}
