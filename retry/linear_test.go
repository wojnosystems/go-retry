package retry_test

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/wojnosystems/go-retry/mocks"
	"github.com/wojnosystems/go-retry/retry"
	"github.com/wojnosystems/go-retry/retryStop"
	"time"
)

const (
	timeUnit = time.Millisecond
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
			retrier retry.Retrier
			mock    *mocks.Callback
		)
		BeforeEach(func() {
			mock = &mocks.Callback{Responses: []error{
				mocks.ErrRetry,
				mocks.ErrRetry,
				mocks.ErrRetry,
				mocks.ErrRetry,
				retryStop.Success,
			}}
			retrier = retry.NewLinear(1*timeUnit, 1.0)
		})
		It("takes the appropriate amount of time", func() {
			elapsed := mocks.DurationElapsed(func() {
				_ = retrier.Retry(ctx, mock.Next())
			})
			Expect(elapsed).Should(BeNumerically(">", 4*timeUnit))
			Expect(elapsed).Should(BeNumerically("<", 14*timeUnit))
		})
	})
})
