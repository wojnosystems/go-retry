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

func loopForever(_ uint64) bool {
	return true
}

func loopNever(_ uint64) bool {
	return false
}

var _ = Describe("Until", func() {

	When("context not expired", func() {
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
		When("the first one succeeds", func() {
			var (
				mock *retryMocks.Callback
			)
			BeforeEach(func() {
				mock = &retryMocks.Callback{
					Responses: []error{
						retryError.StopSuccess,
					},
				}
			})
			It("succeeds", func() {
				err := retryLoop.Until(ctx, mock.Generator(), retryMocks.NeverWaits, loopForever)
				Expect(err).Should(BeNil())
				Expect(mock.TimesRun()).Should(Equal(1))
			})
		})
		When("retries once", func() {
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
				err := retryLoop.Until(ctx, mock.Generator(), retryMocks.NeverWaits, loopForever)
				Expect(err).Should(BeNil())
				Expect(mock.TimesRun()).Should(Equal(2))
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
				err := retryLoop.Until(ctx, mock.Generator(), retryMocks.NeverWaits, loopForever)
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
				err := retryLoop.Until(ctx, mock.Generator(), retryMocks.NeverWaits, loopNever)
				Expect(err).Should(Equal(retryMocks.ErrRetryReason))
				Expect(mock.TimesRun()).Should(Equal(1))
			})
		})
	})
	When("context expired", func() {
		var (
			ctx  context.Context
			mock *retryMocks.Callback
		)
		BeforeEach(func() {
			var cancel context.CancelFunc
			ctx, cancel = context.WithCancel(context.Background())
			cancel()

			mock = &retryMocks.Callback{
				Responses: []error{
					retryMocks.ErrRetry,
					retryMocks.ErrThatCannotBeRetried,
				},
			}
		})
		It("fails without attempting", func() {
			err := retryLoop.Until(ctx, mock.Generator(), retryMocks.NeverWaits, loopForever)
			Expect(err).Should(Equal(context.Canceled))
			Expect(mock.TimesRun()).Should(Equal(0))
		})
	})
})
