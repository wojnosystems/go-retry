package mocks

import (
	"errors"
	"github.com/wojnosystems/go-retry/retryAgain"
)

var (
	ErrRetryReason         = errors.New("forced retry")
	ErrRetry               = retryAgain.Error(ErrRetryReason)
	ErrThatCannotBeRetried = errors.New("un-retryable error")
)
