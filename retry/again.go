package retry

// Again wraps an error and lets the library know it needs to retry
// This was a library decision: most errors should cause the retry system to stop retrying
// Only a certain subset of errors are retryable, usually network-related timeouts
func Again(err error) error {
	return &again{
		wrapped: err,
	}
}

type again struct {
	wrapped error
}

func (a *again) Error() string {
	return a.err().Error()
}

// err Returns the underlying or wrapped error
func (a *again) err() error {
	return a.wrapped
}
