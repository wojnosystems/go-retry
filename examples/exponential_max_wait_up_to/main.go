package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/wojnosystems/go-retry/examples/common"
	"github.com/wojnosystems/go-retry/retry"
	"github.com/wojnosystems/go-retry/retryError"
	"time"
)

func main() {
	tries := 0
	duration := common.TimeThis(func() {
		timer := common.NewTimeSet()

		_ = retry.NewExponentialMaxWaitUpTo(
			10*time.Millisecond,
			1.5,
			10,
			100*time.Millisecond,
		).Retry(context.TODO(), func() (err error) {
			tries++
			fmt.Println(timer.SinceLast())
			return retryError.Again(errors.New("simulated error"))
		})
	})
	fmt.Println("tried", tries, "times taking", duration)
}
