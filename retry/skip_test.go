package retry_test

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/wojnosystems/go-retry/retry"
	"github.com/wojnosystems/go-retry/retryError"
)

var _ = Describe("Skip", func() {
	It("does not call the callback", func() {
		wasCalled := false
		err := retry.Skip.Retry(context.Background(), func() (err error) {
			wasCalled = true
			return retryError.StopSuccess
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(wasCalled).Should(BeFalse())
	})
})
