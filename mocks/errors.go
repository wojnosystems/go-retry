package mocks

import (
	"errors"
	"github.com/wojnosystems/go-retry/retryError"
)

var (
	ErrRetryReason         = errors.New("forced retry")
	ErrRetry               = retryError.Again(ErrRetryReason)
	ErrThatCannotBeRetried = errors.New("un-retryable error")
)
