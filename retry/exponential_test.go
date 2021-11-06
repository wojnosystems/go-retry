package retry_test

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/wojnosystems/go-retry/mocks"
	"github.com/wojnosystems/go-retry/retry"
	"github.com/wojnosystems/go-retry/retryError"
	"time"
)

var _ = Describe("Exponential", func() {
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
				mocks.ErrRetry, // wait 1 * (2)^0 = 1, total 1
				mocks.ErrRetry, // wait 1 * (2)^1 = 2, total 3
				mocks.ErrRetry, // wait 1 * (2)^2 = 4, total 7
				mocks.ErrRetry, // wait 1 * (2)^3 = 8, total 15
				mocks.ErrRetry, // wait 1 * (2)^4 = 16, total 31
				retryError.StopSuccess,
			}}
		})
		When("under retry limit", func() {
			var (
				subject retry.Retrier
			)
			BeforeEach(func() {
				subject = retry.NewExponential(1*timeUnit, 1.0)
			})
			It("takes the appropriate amount of time", func() {
				elapsed := mocks.DurationElapsed(func() {
					_ = subject.Retry(ctx, mock.Next())
				})
				Expect(elapsed).Should(BeNumerically(">", 31*timeUnit))
				Expect(elapsed).Should(BeNumerically("<", 41*timeUnit))
			})
		})
	})
})
