package retryMocks

import (
	"github.com/wojnosystems/go-retry/retryError"
)

// AlwaysSucceeds will always return nil AKA retryError.StopSuccess
func AlwaysSucceeds() error {
	return retryError.StopSuccess
}
