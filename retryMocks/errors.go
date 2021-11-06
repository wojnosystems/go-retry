package retryMocks

import (
	"errors"
	"github.com/wojnosystems/go-retry/retryError"
)

var (
	// ErrRetryReason is a dummy error which could be retried
	ErrRetryReason = errors.New("forced retry")
	// ErrRetry is a dummy error which should be retried, uses ErrRetryReason as the wrapped error
	ErrRetry = retryError.Again(ErrRetryReason)
	// ErrThatCannotBeRetried is an error designated that it should not be used to retry
	ErrThatCannotBeRetried = errors.New("un-retryable error")
)
