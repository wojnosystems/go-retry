package core_test

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/wojnosystems/go-retry/core"
	"github.com/wojnosystems/go-retry/mocks"
	"github.com/wojnosystems/go-retry/retryError"
	"time"
)

func loopForever(_ uint64) bool {
	return true
}

func loopNever(_ uint64) bool {
	return false
}

var _ = Describe("LoopUntil", func() {

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
				mock *mocks.Callback
			)
			BeforeEach(func() {
				mock = &mocks.Callback{
					Responses: []error{
						retryError.StopSuccess,
					},
				}
			})
			It("succeeds", func() {
				err := core.LoopUntil(ctx, mock.Next(), mocks.NeverWaits, loopForever)
				Expect(err).Should(BeNil())
				Expect(mock.TimesRun()).Should(Equal(1))
			})
		})
		When("retries once", func() {
			var (
				mock *mocks.Callback
			)
			BeforeEach(func() {
				mock = &mocks.Callback{
					Responses: []error{
						mocks.ErrRetry,
						retryError.StopSuccess,
					},
				}
			})
			It("succeeds", func() {
				err := core.LoopUntil(ctx, mock.Next(), mocks.NeverWaits, loopForever)
				Expect(err).Should(BeNil())
				Expect(mock.TimesRun()).Should(Equal(2))
			})
		})
		When("unrecoverable error", func() {
			var (
				mock *mocks.Callback
			)
			BeforeEach(func() {
				mock = &mocks.Callback{
					Responses: []error{
						mocks.ErrRetry,
						mocks.ErrThatCannotBeRetried,
					},
				}
			})
			It("fails", func() {
				err := core.LoopUntil(ctx, mock.Next(), mocks.NeverWaits, loopForever)
				Expect(err).Should(Equal(mocks.ErrThatCannotBeRetried))
				Expect(mock.TimesRun()).Should(Equal(2))
			})
		})
		When("retries exceeded", func() {
			var (
				mock *mocks.Callback
			)
			BeforeEach(func() {
				mock = &mocks.Callback{
					Responses: []error{
						mocks.ErrRetry,
						mocks.ErrThatCannotBeRetried,
					},
				}
			})
			It("returns the last retry error", func() {
				err := core.LoopUntil(ctx, mock.Next(), mocks.NeverWaits, loopNever)
				Expect(err).Should(Equal(mocks.ErrRetryReason))
				Expect(mock.TimesRun()).Should(Equal(1))
			})
		})
	})
	When("context expired", func() {
		var (
			ctx  context.Context
			mock *mocks.Callback
		)
		BeforeEach(func() {
			var cancel context.CancelFunc
			ctx, cancel = context.WithCancel(context.Background())
			cancel()

			mock = &mocks.Callback{
				Responses: []error{
					mocks.ErrRetry,
					mocks.ErrThatCannotBeRetried,
				},
			}
		})
		It("fails without attempting", func() {
			err := core.LoopUntil(ctx, mock.Next(), mocks.NeverWaits, loopForever)
			Expect(err).Should(Equal(context.Canceled))
			Expect(mock.TimesRun()).Should(Equal(0))
		})
	})
})
