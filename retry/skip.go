package retry

import (
	"context"
	"github.com/wojnosystems/go-retry/core"
	"github.com/wojnosystems/go-retry/retryError"
)

var Skip = &skip{}

// Skip is a test helping placeholder that will not call the method even once
// You can use this to disable some block of logic guarded with a Retry
type skip struct {
}

func (s *skip) Retry(_ context.Context, _ core.CallbackFunc) (err error) {
	return retryError.StopSuccess
}
