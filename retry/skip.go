package retry

import (
	"context"
	"github.com/wojnosystems/go-retry/retryError"
	"github.com/wojnosystems/go-retry/retryLoop"
)

var Skip = &skip{}

// Skip is a test helping placeholder that will not call the method even once
// You can use this to disable some block of logic guarded with a Retry
type skip struct {
	retryStrategy
}

func (s *skip) Retry(_ context.Context, _ retryLoop.CallbackFunc) (err error) {
	return retryError.StopSuccess
}
