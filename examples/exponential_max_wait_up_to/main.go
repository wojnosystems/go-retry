package main

import (
	"errors"
	"fmt"
	"github.com/wojnosystems/go-retry/examples/common"
	"github.com/wojnosystems/go-retry/retry"
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
		}).Retry(func() (err error) {
			tries++
			fmt.Println(timer.SinceLast())
			return retry.Again(errors.New("simulated error"))
		})

	})
	fmt.Println("tried", tries, "times taking", duration)
}
