package retryMocks

// AlwaysRetries will always return a retryable error
func AlwaysRetries() error {
	return ErrRetryReason
}
