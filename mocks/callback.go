package mocks

// Callback simulates successive calls to a method to be retried
// Responses is a slice of errors to return to each invocation of retry, allowing tests to mock
// the behavior of a thing to retry and validate it
type Callback struct {
	timesRun  int
	Responses []error
}

func (c *Callback) Next() func() error {
	return func() error {
		err := c.Responses[c.timesRun]
		c.timesRun++
		return err
	}
}

func (c *Callback) TimesRun() int {
	return c.timesRun
}
