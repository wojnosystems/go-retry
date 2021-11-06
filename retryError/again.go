package retryError

type AgainWrapper interface {
	// error includes the Error method, forcing AgainWrapper to also be an error type
	error

	// Unwrap is the original error that was wrapped
	Unwrap() error
}

// Again wraps an error and lets the library know it needs to retry the callback
// This was a library decision: most errors should cause the retry system to stop retrying
// Only a certain subset of errors should be retryable, usually network-related timeouts
func Again(err error) AgainWrapper {
	return &again{
		wrapped: err,
	}
}

type again struct {
	wrapped error
}

// Error is the error string of the wrapped Err
func (a *again) Error() string {
	return a.Unwrap().Error()
}

// Unwrap returns the wrapped error given to the Again constructor
func (a *again) Unwrap() error {
	return a.wrapped
}

// IsAgain returns true if this error should be retried by the retry library, false otherwise
func IsAgain(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*again)
	return ok
}
