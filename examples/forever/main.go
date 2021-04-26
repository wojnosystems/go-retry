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

		_ = (&retry.Forever{
			WaitBetweenAttempts: 10 * time.Millisecond,
		}).Retry(func() (err error) {
			fmt.Println(timer.SinceLast())
			if tries < 10 {
				tries++
				return retry.Again(errors.New("simulated error"))
			}
			return retry.Success
		})
	})
	fmt.Println("tried", tries, "times taking", duration)
}
