package retryMocks

// AlwaysFails is a call-back that will always return a non-retryable error
func AlwaysFails() error {
	return ErrThatCannotBeRetried
}
