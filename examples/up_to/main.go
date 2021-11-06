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
	timer := common.NewTimeSet()
	duration := common.TimeThis(func() {

		_ = retry.NewUpTo(
			10*time.Millisecond,
			10,
		).Retry(context.TODO(), func() (err error) {
			fmt.Println(timer.SinceLast())
			tries++
			return retryError.Again(errors.New("simulated error"))
		})
	})
	fmt.Println("tried", tries, "times taking", duration)
}
