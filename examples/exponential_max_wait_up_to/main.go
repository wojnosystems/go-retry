package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/wojnosystems/go-retry/examples/common"
	"github.com/wojnosystems/go-retry/retry"
	"github.com/wojnosystems/go-retry/retryAgain"
	"time"
)

func main() {
	tries := 0
	duration := common.TimeThis(func() {

		timer := common.NewTimeSet()

		_ = (&retry.ExponentialMaxWaitUpTo{
			InitialWaitBetweenAttempts: 10 * time.Millisecond,
			GrowthFactor:               1.5,
			MaxAttempts:                10,
			MaxWaitBetweenAttempts:     100 * time.Millisecond,
		}).Retry(context.TODO(), func() (err error) {
			tries++
			fmt.Println(timer.SinceLast())
			return retryAgain.Error(errors.New("simulated error"))
		})

	})
	fmt.Println("tried", tries, "times taking", duration)
}
