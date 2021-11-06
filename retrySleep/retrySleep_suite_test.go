package retrySleep_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRetrySleep(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RetrySleep Suite")
}
