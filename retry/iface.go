package retry

import "github.com/wojnosystems/go-retry/core"

type Retrier interface {
	Retry(cb core.CallbackFunc) (err error)
}
