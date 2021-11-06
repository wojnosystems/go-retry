package retryLoop_test

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/wojnosystems/go-retry/retryError"
	"github.com/wojnosystems/go-retry/retryLoop"
	"github.com/wojnosystems/go-retry/retryMocks"
	"time"
)

const (
	noRetries   = 0
	manyRetries = 1000
)

var _ = Describe("UpTo", func() {
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
	When("retry attempts available", func() {
		var (
			mock *retryMocks.Callback
		)
		BeforeEach(func() {
			mock = &retryMocks.Callback{
				Responses: []error{
					retryMocks.ErrRetry,
					retryError.StopSuccess,
				},
			}
		})
		It("succeeds", func() {
			err := retryLoop.UpTo(ctx, mock.Generator(), retryMocks.NeverWaits, manyRetries)
			Expect(err).Should(BeNil())
			Expect(mock.TimesRun()).Should(Equal(2))
		})
	})
	When("no retry attempts available", func() {
		var (
			mock *retryMocks.Callback
		)
		BeforeEach(func() {
			mock = &retryMocks.Callback{
				Responses: []error{
					retryMocks.ErrRetry,
					retryError.StopSuccess,
				},
			}
		})
		It("stops retrying", func() {
			err := retryLoop.UpTo(ctx, mock.Generator(), retryMocks.NeverWaits, noRetries)
			Expect(err).Should(Equal(retryMocks.ErrRetryReason))
			Expect(mock.TimesRun()).Should(Equal(1))
		})
	})
	When("unrecoverable error", func() {
		var (
			mock *retryMocks.Callback
		)
		BeforeEach(func() {
			mock = &retryMocks.Callback{
				Responses: []error{
					retryMocks.ErrRetry,
					retryMocks.ErrThatCannotBeRetried,
				},
			}
		})
		It("fails", func() {
			err := retryLoop.UpTo(ctx, mock.Generator(), retryMocks.NeverWaits, manyRetries)
			Expect(err).Should(Equal(retryMocks.ErrThatCannotBeRetried))
			Expect(mock.TimesRun()).Should(Equal(2))
		})
	})
	When("retries exceeded", func() {
		var (
			mock *retryMocks.Callback
		)
		BeforeEach(func() {
			mock = &retryMocks.Callback{
				Responses: []error{
					retryMocks.ErrRetry,
					retryMocks.ErrThatCannotBeRetried,
				},
			}
		})
		It("returns the last retry error", func() {
			err := retryLoop.UpTo(ctx, mock.Generator(), retryMocks.NeverWaits, noRetries)
			Expect(err).Should(Equal(retryMocks.ErrRetryReason))
			Expect(mock.TimesRun()).Should(Equal(1))
		})
	})
})
