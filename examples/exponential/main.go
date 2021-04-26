package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/wojnosystems/go-retry/examples/common"
	"github.com/wojnosystems/go-retry/retry"
	"github.com/wojnosystems/go-retry/retryAgain"
	"github.com/wojnosystems/go-retry/retryStop"
	"time"
)

func main() {
	tries := 0
	duration := common.TimeThis(func() {

		timer := common.NewTimeSet()

		_ = (&retry.Exponential{
			InitialWaitBetweenAttempts: 10 * time.Millisecond,
			GrowthFactor:               1.0,
		}).Retry(context.TODO(), func() (err error) {
			fmt.Println(timer.SinceLast())
			if tries < 8 {
				tries++
				return retryAgain.Error(errors.New("simulated error"))
			}
			return retryStop.Success
		})

	})
	fmt.Println("tried", tries, "times taking", duration)
}
