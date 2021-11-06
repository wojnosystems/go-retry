package mocks

import (
	"github.com/wojnosystems/go-retry/retryError"
)

func AlwaysSucceeds() func() error {
	return func() error {
		return retryError.StopSuccess
	}
}
