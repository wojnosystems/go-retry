package core

// CallbackFunc is called each time a retryable attempt needs to be made
// return nil AKA retryStop.Success to stop retrying
// return retryAgain.Error(err) to retry. If no more attempts can be made, the Error will be returned to the caller
// return any other error to stop retrying and return the error immediately
type CallbackFunc func() (err error)

// DelayBetweenAttemptsFunc is called after a retryable error is received and additional
// retry attempts remain. If there are no retry attempts, this method will not be called.
// iterator is the current number of delays that have been called (starts at 0).
// This will count up monotonically until it reaches math.MaxUInt64, at which point, it will stop counting count and return math.MaxUInt64.
// iterator is intended to allow developers to generate delays based on number of times the request failed before.
type DelayBetweenAttemptsFunc func(iterator uint64)
