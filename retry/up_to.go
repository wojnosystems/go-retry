package retry

import "time"

// UpTo retries up to MaxAttempts and WaitBetweenAttempts duration between each retryable error
type UpTo struct {
	// WaitBetweenAttempts
	WaitBetweenAttempts time.Duration

	// MaxAttempts is how many failed tries to attempt before returning an error and giving up
	MaxAttempts uint
}

func (c *UpTo) Retry(cb func() (err error)) (err error) {
	return loopUpTo(cb, func(i uint64) {
		time.Sleep(c.WaitBetweenAttempts)
	}, uint64(c.MaxAttempts))
}

func loopUpTo(cb func() error, wait func(i uint64), maxAttempts uint64) (err error) {
	i := uint64(0)
	for {
		err = cb()
		if err == Success {
			return
		}
		if v, ok := err.(*again); !ok {
			return err
		} else {
			i++
			if i < maxAttempts {
				wait(i - 1)
			} else {
				return v.err()
			}
		}
	}
}
