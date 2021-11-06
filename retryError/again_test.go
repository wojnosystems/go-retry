package retryError

import (
	"errors"
	. "github.com/onsi/gomega"
	"testing"
)

var errFake = errors.New("fake")

func TestAgain_Err(t *testing.T) {
	g := NewWithT(t)
	err := Again(errFake)
	g.Expect(err).Should(HaveOccurred())
	g.Expect(err.Unwrap()).Should(Equal(errFake))
}

func TestAgain_Error(t *testing.T) {
	g := NewWithT(t)
	err := Again(errFake)
	g.Expect(err.Error()).Should(Equal(errFake.Error()))
}

func TestAgain_IsAgain(t *testing.T) {
	cases := map[string]struct {
		input    error
		expected bool
	}{
		"nil": {
			input: nil,
		},
		"retryable": {
			input:    Again(errFake),
			expected: true,
		},
		"not retryable": {
			input: errFake,
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			g := NewWithT(t)
			actual := IsAgain(c.input)
			g.Expect(actual).Should(Equal(c.expected))
		})
	}
}
