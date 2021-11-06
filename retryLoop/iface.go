package retryLoop

// CallbackFunc is called each time a retryable attempt needs to be made
// return nil AKA retryStop.Success to stop retrying
// return retryAgain.Error(err) to retry. If no more attempts can be made, the Error will be returned to the caller
// return any other error to stop retrying and return the error immediately
type CallbackFunc func() (err error)

// WaitBetweenAttemptsFunc is called after a retryable error is received and there are additional retry attempts
// permitted. If there are no retry attempts, this method will not be called.
// timesWaited is the current number of wait that have been requested (starts at 0).
// This will count up monotonically until it reaches math.MaxUInt64, at which point, it will stop counting up and will
// always return the maximum uint64 value. timesWaited is intended to allow developers to generate delays based on
// number of times the request failed.
type WaitBetweenAttemptsFunc func(timesWaited uint64)

// ShouldContinueLoopingFunc should return true if another retry should be attempted, false to stop
// This method is only called if an error wrapped in a retryError.Again is returned
// timesAttempted will return the number of times the call-back has been attempted (starts at 1) and will count up until it reaches the maximum uint64 size, at which point it will stop counting, but continue calling your method
type ShouldContinueLoopingFunc func(timesAttempted uint64) bool
