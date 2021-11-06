package mocks

import "github.com/wojnosystems/go-retry/retryStop"

func AlwaysSucceeds() func() error {
	return func() error {
		return retryStop.Success
	}
}
