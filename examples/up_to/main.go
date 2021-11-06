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
	// you can define your retry strategy and use it wherever you need
	// you can use the struct or retry.NewUpTo constructor, whichever you find cleaner
	retryStrategy := &retry.UpTo{
		// we will wait 10ms between every attempt that failed but could be retried
		WaitBetweenAttempts: 10 * time.Millisecond,
		// we will only attempt the callback 10 times before returning the last retryable error
		MaxAttempts: 10,
	}

	tries := 0
	timer := common.NewTimeSet()
	var err error
	duration := common.TimeThis(func() {
		err = retryStrategy.Retry(context.TODO(), func() (err error) {
			fmt.Println(timer.SinceLast())
			tries++
			return retryError.Again(errors.New("simulated error"))
		})
	})
	fmt.Println("tried", tries, "times taking", duration)
	fmt.Println("should get 'simulated error': ", err.Error())
}
