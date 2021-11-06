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

const (
	noRetries   = 0
	manyRetries = 1000
)

var _ = Describe("LoopUpTo", func() {
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
			err := core.LoopUpTo(ctx, mock.Next(), mocks.NeverWaits, manyRetries)
			Expect(err).Should(BeNil())
			Expect(mock.TimesRun()).Should(Equal(2))
		})
	})
	When("no retry attempts available", func() {
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
		It("stops retrying", func() {
			err := core.LoopUpTo(ctx, mock.Next(), mocks.NeverWaits, noRetries)
			Expect(err).Should(Equal(mocks.ErrRetryReason))
			Expect(mock.TimesRun()).Should(Equal(1))
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
			err := core.LoopUpTo(ctx, mock.Next(), mocks.NeverWaits, manyRetries)
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
			err := core.LoopUpTo(ctx, mock.Next(), mocks.NeverWaits, noRetries)
			Expect(err).Should(Equal(mocks.ErrRetryReason))
			Expect(mock.TimesRun()).Should(Equal(1))
		})
	})
})
