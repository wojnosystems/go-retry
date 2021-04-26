package main

import (
	"fmt"
	"github.com/wojnosystems/go-retry/examples/common"
	"github.com/wojnosystems/go-retry/retry"
	"net"
	"time"
)

func main() {
	dialer := &retry.ExponentialMaxWaitUpTo{
		InitialWaitBetweenAttempts: 50 * time.Millisecond,
		GrowthFactor:               1.0,
		MaxAttempts:                8,
		MaxWaitBetweenAttempts:     500 * time.Millisecond,
	}

	timer := common.NewTimeSet()

	totalTime := common.TimeThis(func() {
		err := dialer.Retry(func() error {
			fmt.Println("dialing", timer.SinceLast())
			socket, dialErr := net.Dial("tcp", "localhost:9999")
			if dialErr != nil {
				// all dialErrs are retried
				return retry.Again(dialErr)
			}

			// Write errors are NOT retried
			_, writeErr := socket.Write([]byte("some payload"))

			// if writeErr is nil, success!
			// if writeErr is not wrapped in retry.Again, retry will stop retrying and return the
			// error to the caller
			return writeErr
		})

		// Outputs the Dial error because we ran out of retries
		fmt.Println(err)
	})

	fmt.Println("total time", totalTime)
}
