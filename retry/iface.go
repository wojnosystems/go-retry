package retry

import (
	"context"
	"github.com/wojnosystems/go-retry/retryLoop"
)

type Retrier interface {
	Retry(ctx context.Context, cb retryLoop.CallbackFunc) (err error)
}
