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

var _ = Describe("Forever", func() {
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
				retryMocks.ErrRetry, // wait 1, total 1
				retryMocks.ErrRetry, // wait 1, total 2
				retryMocks.ErrRetry, // wait 1, total 3
				retryMocks.ErrRetry, // wait 1, total 4
				retryMocks.ErrRetry, // wait 1, total 5
				retryError.StopSuccess,
			}}
		})
		When("under retry limit", func() {
			var (
				subject retry.Retrier
			)
			BeforeEach(func() {
				subject = retry.NewForever(1 * timeUnit)
			})
			It("takes the appropriate amount of time", func() {
				elapsed := retryMocks.DurationElapsed(func() {
					_ = subject.Retry(ctx, mock.Generator())
				})
				Expect(elapsed).Should(BeNumerically(">", 5*timeUnit))
				Expect(elapsed).Should(BeNumerically("<", 10*timeUnit))
			})
		})
	})
})
