package main

import (
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

		_ = (&retry.UpTo{
			WaitBetweenAttempts: 10 * time.Millisecond,
			MaxAttempts:         10,
		}).Retry(func() (err error) {
			fmt.Println(timer.SinceLast())
			tries++
			return retryAgain.Error(errors.New("simulated error"))
		})

	})
	fmt.Println("tried", tries, "times taking", duration)
}
