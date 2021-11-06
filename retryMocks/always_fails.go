package retryMocks

// AlwaysFails will always return a non-retryable error
func AlwaysFails() error {
	return ErrThatCannotBeRetried
}
