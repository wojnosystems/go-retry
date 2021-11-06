package mocks

func AlwaysRetries() func() error {
	return func() error {
		return ErrRetryReason
	}
}
