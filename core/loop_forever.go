package core

import (
	"github.com/wojnosystems/go-retry/retryAgain"
	"github.com/wojnosystems/go-retry/retryStop"
	"math"
)

// LoopForever will continuously call the callback (cb) until it succeeds or
// returns a non-retryable error
func LoopForever(cb CallbackFunc, wait DelayBetweenAttemptsFunc) (err error) {
	i := uint64(0)
	for {
		err = cb()
		if err == retryStop.Success {
			return
		}
		if !retryAgain.IsAgain(err) {
			return err
		}
		wait(i)
		if i < math.MaxUint64 {
			i++
		}
	}
}
