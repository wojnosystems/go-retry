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

var _ = Describe("LoopForever", func() {
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
	When("multiple retries", func() {
		var (
			mock *mocks.Callback
		)
		BeforeEach(func() {
			mock = &mocks.Callback{
				Responses: []error{
					mocks.ErrRetry,
					mocks.ErrRetry,
					mocks.ErrRetry,
					mocks.ErrRetry,
					mocks.ErrRetry,
					mocks.ErrRetry,
					retryError.StopSuccess,
				},
			}
		})
		It("does not exhaust retries", func() {
			err := core.LoopForever(ctx, mock.Next(), mocks.NeverWaits)
			Expect(err).Should(BeNil())
			Expect(mock.TimesRun()).Should(Equal(7))
		})
	})
})
