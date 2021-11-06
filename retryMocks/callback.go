package retryMocks

// Callback simulates successive calls to a method to be retried
// Responses is a slice of errors to return to each invocation of retry, allowing tests to mock
// the behavior of a thing to retry and validate it
type Callback struct {
	timesRun  int
	Responses []error
}

// Generator returns a handle to the callback method that will be called by the Retrier to return each of the
// Responses in order. If you try to call it more times than there are Responses, it will panic
func (c *Callback) Generator() func() error {
	return func() error {
		err := c.Responses[c.timesRun]
		c.timesRun++
		return err
	}
}

// TimesRun gets the number of times Generator's returned function was called
func (c *Callback) TimesRun() int {
	return c.timesRun
}
