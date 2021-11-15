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
			subject *retry.Linear
			mock    *retryMocks.Callback
		)
		BeforeEach(func() {
			mock = &retryMocks.Callback{Responses: []error{
				retryMocks.ErrRetry,
				retryMocks.ErrRetry,
				retryMocks.ErrRetry,
				retryMocks.ErrRetry,
				retryError.StopSuccess,
			}}
			subject = retry.NewLinear(1*timeUnit, 1.0)
		})
		It("takes the appropriate amount of time", func() {
			elapsed := retryMocks.DurationElapsed(func() {
				_ = subject.Retry(ctx, mock.Generator())
			})
			Expect(elapsed).Should(BeNumerically(">", 4*timeUnit))
			Expect(elapsed).Should(BeNumerically("<", 14*timeUnit))
		})
	})
})
