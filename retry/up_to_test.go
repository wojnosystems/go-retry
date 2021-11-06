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

var _ = Describe("UpTo", func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
		mock   *retryMocks.Callback
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	})

	AfterEach(func() {
		cancel()
	})

	When("succeeds the first time", func() {
		It("does not wait", func() {
			subject := retry.NewUpTo(1*time.Hour, 1000)
			err := subject.Retry(ctx, retryMocks.AlwaysSucceeds)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	When("succeeds after 5 times", func() {
		BeforeEach(func() {
			mock = &retryMocks.Callback{
				Responses: []error{
					retryError.Again(retryMocks.ErrRetryReason),
					retryError.Again(retryMocks.ErrRetryReason),
					retryError.Again(retryMocks.ErrRetryReason),
					retryError.Again(retryMocks.ErrRetryReason),
					retryError.Again(retryMocks.ErrRetryReason),
					retryError.StopSuccess,
				},
			}
		})
		It("succeeds", func() {
			subject := retry.NewUpTo(0, 10000)
			err := subject.Retry(ctx, mock.Generator())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mock.TimesRun()).Should(Equal(6))
		})
		It("takes at least 4 milliseconds", func() {
			subject := &retry.UpTo{MaxAttempts: 10000, WaitBetweenAttempts: 1 * timeUnit}
			elapsed := retryMocks.DurationElapsed(func() {
				_ = subject.Retry(ctx, mock.Generator())
			})
			Expect(elapsed).Should(BeNumerically(">=", 4*timeUnit))
		})
	})
	When("retries exhausted", func() {
		It("fails", func() {
			subject := retry.NewUpTo(0, 1)
			err := subject.Retry(ctx, retryMocks.AlwaysRetries)
			Expect(err).Should(Equal(retryMocks.ErrRetryReason))
		})
	})
	When("error is not retryable", func() {
		It("fails", func() {
			subject := retry.NewUpTo(0, 1)
			err := subject.Retry(ctx, retryMocks.AlwaysFails)
			Expect(err).Should(Equal(retryMocks.ErrThatCannotBeRetried))
		})
	})
})
