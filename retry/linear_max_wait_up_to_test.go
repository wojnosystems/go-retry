package retry_test

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/wojnosystems/go-retry/retry"
	"github.com/wojnosystems/go-retry/retryError"
	"github.com/wojnosystems/go-retry/retryMocks"
	"time"
)

var _ = Describe("LinearRetry", func() {
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
			mock *retryMocks.Callback
		)
		BeforeEach(func() {
			mock = &retryMocks.Callback{Responses: []error{
				retryMocks.ErrRetry, // 1 + (1*1*0) = 1, total 1
				retryMocks.ErrRetry, // 1 + (1*1*1) = 2, total 3
				retryMocks.ErrRetry, // 1 + (1*1*2) = 3, total 6
				retryMocks.ErrRetry, // 1 + (1*1*3) = 4, total 10
				retryMocks.ErrRetry, // 1 + (1*1*4) = 5, total 15
				retryMocks.ErrRetry, // 1 + (1*1*5) = 6 (5), total 21 (20) (5 is max wait time)
				retryMocks.ErrRetry, // 1 + (1*1*6) = 7 (5), total 28 (25)
				retryMocks.ErrRetry, // 1 + (1*1*6) = 8 (5), total 36 (30)
				retryMocks.ErrRetry, // 1 + (1*1*6) = 9 (5), total 45 (35)
				retryError.StopSuccess,
			}}
		})
		When("under retry limit", func() {
			var (
				subject *retry.LinearMaxWaitUpTo
			)
			BeforeEach(func() {
				subject = retry.NewLinearMaxWaitUpTo(1*timeUnit, 1.0, 10, 5*timeUnit)
			})
			It("takes the appropriate amount of time", func() {
				elapsed := retryMocks.DurationElapsed(func() {
					_ = subject.Retry(ctx, mock.Generator())
				})
				Expect(elapsed).Should(BeNumerically(">", 35*timeUnit))
				Expect(elapsed).Should(BeNumerically("<", 45*timeUnit))
			})
		})
		When("over retry limit", func() {
			var (
				subject *retry.LinearMaxWaitUpTo
			)
			BeforeEach(func() {
				subject = retry.NewLinearMaxWaitUpTo(1*timeUnit, 1.0, 4, 5*timeUnit)
			})
			It("takes the appropriate amount of time", func() {
				elapsed := retryMocks.DurationElapsed(func() {
					_ = subject.Retry(ctx, mock.Generator())
				})
				Expect(elapsed).Should(BeNumerically(">", 6*timeUnit))
				Expect(elapsed).Should(BeNumerically("<", 16*timeUnit))
			})
		})
	})
})
