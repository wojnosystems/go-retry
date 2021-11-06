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

		_ = retry.NewExponential(
			10*time.Millisecond,
			1.0,
		).Retry(context.TODO(), func() (err error) {
			fmt.Println(timer.SinceLast())
			if tries < 8 {
				tries++
				return retryError.Again(errors.New("simulated error"))
			}
			return nil
		})
	})
	fmt.Println("tried", tries, "times taking", duration)
}
