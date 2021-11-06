package retry_test

import (
	"context"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/wojnosystems/go-retry/retry"
	"github.com/wojnosystems/go-retry/retryAgain"
	"github.com/wojnosystems/go-retry/retryStop"
)

var _ = Describe("Never", func() {
	When("success", func() {
		It("calls back exactly once", func() {
			wasCalled := false
			err := retry.Never.Retry(context.Background(), func() (err error) {
				wasCalled = true
				return retryStop.Success
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(wasCalled).Should(BeTrue())
		})
	})
	When("fail", func() {
		It("calls back exactly once", func() {
			expectedErr := errors.New("intentionally fails")
			wasCalled := false
			err := retry.Never.Retry(context.Background(), func() (err error) {
				wasCalled = true
				return retryAgain.Error(expectedErr)
			})
			Expect(err).Should(HaveOccurred())
			Expect(wasCalled).Should(BeTrue())
		})
	})
})
