package retryAgain

type Wrapper interface {
	Error() string
	Err() error
}
