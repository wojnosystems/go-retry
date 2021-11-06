package retryLoop

import (
	"context"
	"github.com/wojnosystems/go-retry/retryError"
	"math"
)

// Until will continuously call the callback until shouldContinueLooping returns false.
// Wait is called after each attempt if more attempts should be made. It is expected that "wait" will call sleep or use time.After
// to ensure that the retry waits the appropriate amount of time before trying again.
// Should ctx be canceled or expire, the current callback will be allowed to finish, but no additional attempts or
// waits will be allowed to occur, or complete. That way, you should never over-wait the ctx deadline by a significant amount.
// Wait may still be called after a context expires, wait is expected to take the context into account and only sleep
// until the deadline expires or the retry wait duration expires, whichever occurs first.
// This method is the base for all retry logic. Both Forever and UpTo are intended to depend on this.
func Until(ctx context.Context,
	callback CallbackFunc,
	wait WaitBetweenAttemptsFunc,
	shouldContinueLooping ShouldContinueLoopingFunc,
) (err error) {

	timesAttempted := uint64(0)
	for {
		// Check if context is done, if not, continue
		select {
		case <-ctx.Done():
			// ctx expired or was cancelled, we're done
			return ctx.Err()
		default:
			// fall-through, ctx is not done
		}
		// call the callback, record the response
		err = callback()
		if err == retryError.StopSuccess {
			// attempt succeeded, no need to wait or try again
			return
		}
		if v, ok := err.(retryError.AgainWrapper); !ok {
			// error was no retryable, stop retrying without waiting
			return err
		} else {
			// error was retryable
			if timesAttempted != math.MaxUint64 {
				// only count up if that's possible, avoid overflow
				timesAttempted++
			}
			if shouldContinueLooping(timesAttempted) {
				// we should continue looping, so wait before trying again
				wait(timesAttempted - 1)
			} else {
				// we should not loop again, just return the last error we got, without the retryAgain wrapper
				return v.Unwrap()
			}
		}
	}
}
