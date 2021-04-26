package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/wojnosystems/go-retry/examples/common"
	"github.com/wojnosystems/go-retry/retry"
	"github.com/wojnosystems/go-retry/retryAgain"
	"net/http"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	dialer := &retry.ExponentialMaxWaitUpTo{
		InitialWaitBetweenAttempts: 50 * time.Millisecond,
		GrowthFactor:               1.0,
		MaxAttempts:                15,
		MaxWaitBetweenAttempts:     500 * time.Millisecond,
	}

	timer := common.NewTimeSet()

	totalTime := common.TimeThis(func() {
		err := dialer.Retry(ctx, func() error {
			fmt.Println("getting", timer.SinceLast())
			req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/non-existent", nil)
			req = req.WithContext(ctx)
			_, getErr := http.DefaultClient.Do(req)
			if getErr != nil {
				getErr = errors.Unwrap(getErr)
				if getErr != context.DeadlineExceeded {
					getErr = retryAgain.Error(getErr)
				}
				return getErr
			}
			return nil
		})

		// Outputs the http error because we ran out of retries
		fmt.Println(err)
	})

	fmt.Println("total time", totalTime)
}
