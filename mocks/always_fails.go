package mocks

func AlwaysFails() func() error {
	return func() error {
		return ErrThatCannotBeRetried
	}
}
