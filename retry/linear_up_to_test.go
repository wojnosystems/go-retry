package retry_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/wojnosystems/go-retry/retryError"
	"time"
)

import (
	"context"
	. "github.com/onsi/gomega"
	"github.com/wojnosystems/go-retry/mocks"
	"github.com/wojnosystems/go-retry/retry"
)

var _ = Describe("LinearUpTo", func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	})
	AfterEach(func() {
		cancel()
	})

	When("multiple failures", func() {
		var (
			mock *mocks.Callback
		)
		BeforeEach(func() {
			mock = &mocks.Callback{Responses: []error{
				mocks.ErrRetry, // 1 + (1*1*0) = 1, total 1
				mocks.ErrRetry, // 1 + (1*1*1) = 2, total 3
				mocks.ErrRetry, // 1 + (1*1*2) = 3, total 6
				mocks.ErrRetry, // 1 + (1*1*3) = 4, total 10
				mocks.ErrRetry, // 1 + (1*1*4) = 5, total 15
				mocks.ErrRetry, // 1 + (1*1*5) = 6, total 21
				retryError.StopSuccess,
			}}
		})
		When("under retry limit", func() {
			var (
				subject retry.Retrier
			)
			BeforeEach(func() {
				subject = retry.NewLinearUpTo(1*timeUnit, 1.0, 10)
			})
			It("takes the appropriate amount of time", func() {
				elapsed := mocks.DurationElapsed(func() {
					_ = subject.Retry(ctx, mock.Next())
				})
				Expect(elapsed).Should(BeNumerically(">", 21*timeUnit))
				Expect(elapsed).Should(BeNumerically("<", 31*timeUnit))
			})
		})
		When("over retry limit", func() {
			var (
				subject retry.Retrier
			)
			BeforeEach(func() {
				subject = retry.NewLinearUpTo(1*timeUnit, 1.0, 4)
			})
			It("takes the appropriate amount of time", func() {
				elapsed := mocks.DurationElapsed(func() {
					_ = subject.Retry(ctx, mock.Next())
				})
				Expect(elapsed).Should(BeNumerically(">", 6*timeUnit))
				Expect(elapsed).Should(BeNumerically("<", 12*timeUnit))
			})
		})
	})
})
