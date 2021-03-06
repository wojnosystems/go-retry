package retrySleep_test

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/wojnosystems/go-retry/retryMocks"
	"github.com/wojnosystems/go-retry/retrySleep"
	"time"
)

const (
	aVeryLongTime  = 24 * time.Hour
	aShortTime     = 1 * time.Minute
	aVeryShortTime = 1 * time.Millisecond
)

var _ = Describe("Sleep", func() {
	When("wait is longer than context", func() {
		var (
			ctx context.Context
		)
		BeforeEach(func() {
			var cancel context.CancelFunc
			ctx, cancel = context.WithCancel(context.Background())
			cancel()
		})
		It("does not wait", func() {
			elapsed := retryMocks.DurationElapsed(func() {
				retrySleep.WithContext(ctx, aVeryLongTime)
			})
			Expect(elapsed).Should(BeNumerically("<", aShortTime))
		})
	})
	When("context is longer than wait", func() {
		var (
			ctx    context.Context
			cancel context.CancelFunc
		)
		BeforeEach(func() {
			ctx, cancel = context.WithCancel(context.Background())
		})
		AfterEach(func() {
			cancel()
		})
		It("does not wait", func() {
			elapsed := retryMocks.DurationElapsed(func() {
				retrySleep.WithContext(ctx, aVeryShortTime)
			})
			Expect(elapsed).Should(BeNumerically("~", aVeryShortTime, 1*time.Millisecond))
		})
	})
})
