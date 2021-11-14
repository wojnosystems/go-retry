package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/wojnosystems/go-retry/retry"
	"github.com/wojnosystems/go-retry/retryError"
	"github.com/wojnosystems/go-retry/retryLoop"
	"time"
)

type retryStrategy interface {
	Retry(ctx context.Context, cb retryLoop.CallbackFunc) (err error)
}

func main() {
	normal := &retry.LinearUpTo{
		InitialWaitBetweenAttempts: 100 * time.Millisecond,
		GrowthFactor:               1,
		MaxAttempts:                5,
	}

	var strategy retryStrategy
	strategy = normal

	_ = strategy.Retry(context.TODO(), func() (err error) {
		fmt.Println("normal")
		return retryError.Again(errors.New("some error"))
	})

	strategy = retry.Never

	_ = strategy.Retry(context.TODO(), func() (err error) {
		fmt.Println("NEVER")
		return retryError.Again(errors.New("some error"))
	})

	strategy = normal
	_ = strategy.Retry(context.TODO(), func() (err error) {
		fmt.Println("normal")
		return retryError.Again(errors.New("some error"))
	})
}
