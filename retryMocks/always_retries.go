package retryMocks

// AlwaysRetries is a call-back that will always return a retryable error
func AlwaysRetries() error {
	return ErrRetryReason
}
