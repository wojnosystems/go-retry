package retryAgain

// Error wraps an error and lets the library know it needs to retry
// This was a library decision: most errors should cause the retry system to stop retrying
// Only a certain subset of errors are retryable, usually network-related timeouts
func Error(err error) Wrapper {
	return &again{
		wrapped: err,
	}
}

type again struct {
	wrapped error
}

func (a *again) Error() string {
	return a.Err().Error()
}

// Err Returns the underlying or wrapped error
func (a *again) Err() error {
	return a.wrapped
}

func IsAgain(err error) bool {
	_, ok := err.(*again)
	return ok
}
