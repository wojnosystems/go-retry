package retry_test

import (
	"context"
	"github.com/wojnosystems/go-retry/retryLoop"
)

type retryStrategy interface {
	Retry(ctx context.Context, cb retryLoop.CallbackFunc) (err error)
}
