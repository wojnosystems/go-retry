package retry

import "time"

// Forever will retry forever until no errors are reported
type Forever struct {
	WaitBetweenAttempts time.Duration
}

func (c *Forever) Retry(cb func() (err error)) (err error) {
	return loopForever(cb, func() {
		time.Sleep(c.WaitBetweenAttempts)
	})
}

func loopForever(cb func() error, wait func()) (err error) {
	for {
		err = cb()
		if err == Success {
			return
		}
		if _, ok := err.(*again); !ok {
			return err
		} else {
			wait()
		}
	}
}
