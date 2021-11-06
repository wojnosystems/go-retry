package retry

import (
	"github.com/onsi/gomega"
	"testing"
	"time"
)

func TestMinDuration(t *testing.T) {
	cases := map[string]struct {
		a, b     time.Duration
		expected time.Duration
	}{
		"a > b": {
			a:        20 * time.Millisecond,
			b:        10 * time.Millisecond,
			expected: 10 * time.Millisecond,
		},
		"b > a": {
			a:        10 * time.Millisecond,
			b:        20 * time.Millisecond,
			expected: 10 * time.Millisecond,
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			g := gomega.NewWithT(t)
			actual := minDuration(c.a, c.b)
			g.Expect(actual).Should(gomega.Equal(c.expected))
		})
	}
}
