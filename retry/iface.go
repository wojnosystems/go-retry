package retry

import (
	"context"
	"github.com/wojnosystems/go-retry/core"
)

type Retrier interface {
	Retry(ctx context.Context, cb core.CallbackFunc) (err error)
}
