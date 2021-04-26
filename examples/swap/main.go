package main

import (
	"errors"
	"fmt"
	"github.com/wojnosystems/go-retry/retry"
	"github.com/wojnosystems/go-retry/retryAgain"
	"time"
)

func main() {
	normal := &retry.LinearUpTo{
		InitialWaitBetweenAttempts: 100 * time.Millisecond,
		GrowthFactor:               1,
		MaxAttempts:                5,
	}

	var strategy retry.Retrier
	strategy = normal

	_ = strategy.Retry(func() (err error) {
		fmt.Println("normal")
		return retryAgain.Error(errors.New("some error"))
	})

	strategy = retry.Never

	_ = strategy.Retry(func() (err error) {
		fmt.Println("NEVER")
		return retryAgain.Error(errors.New("some error"))
	})

	strategy = normal
	_ = strategy.Retry(func() (err error) {
		fmt.Println("normal")
		return retryAgain.Error(errors.New("some error"))
	})
}
