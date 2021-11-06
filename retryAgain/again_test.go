package retryAgain

import (
	"errors"
	. "github.com/onsi/gomega"
	"testing"
)

var errFake = errors.New("fake")

func TestAgain_Err(t *testing.T) {
	g := NewWithT(t)
	err := Error(errFake)
	g.Expect(err).Should(HaveOccurred())
	g.Expect(err.Err()).Should(Equal(errFake))
}

func TestAgain_Error(t *testing.T) {
	g := NewWithT(t)
	err := Error(errFake)
	g.Expect(err.Error()).Should(Equal(errFake.Error()))
}

func TestAgain_IsAgain(t *testing.T) {
	g := NewWithT(t)
	err := Error(errFake)
	g.Expect(IsAgain(err)).Should(BeTrue())
}
