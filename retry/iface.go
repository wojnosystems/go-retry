package retry

var Success = error(nil)

type Retrier interface {
	Retry(func() (err error)) (err error)
}
