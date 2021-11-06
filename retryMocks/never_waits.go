package retryMocks

// NeverWaits is a mock wait method to ensure tests induce no delay between attempts
func NeverWaits(_ uint64) {
	// intentionally blank
}
