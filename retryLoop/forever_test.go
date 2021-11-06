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
	When("multiple retries", func() {
		var (
			mock *retryMocks.Callback
		)
		BeforeEach(func() {
			mock = &retryMocks.Callback{
				Responses: []error{
					retryMocks.ErrRetry,
					retryMocks.ErrRetry,
					retryMocks.ErrRetry,
					retryMocks.ErrRetry,
					retryMocks.ErrRetry,
					retryMocks.ErrRetry,
					retryError.StopSuccess,
				},
			}
		})
		It("does not exhaust retries", func() {
			err := retryLoop.Forever(ctx, mock.Generator(), retryMocks.NeverWaits)
			Expect(err).Should(BeNil())
			Expect(mock.TimesRun()).Should(Equal(7))
		})
	})
})
