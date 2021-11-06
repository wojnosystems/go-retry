package retry

import (
	"github.com/onsi/gomega"
	"strconv"
	"testing"
	"time"
)

func Test_exponentialSleepTime(t *testing.T) {
	cases := []struct {
		initial      time.Duration
		growthFactor float64
		iterations   uint64
		expected     time.Duration
	}{
		{
			initial:      1 * timeUnit,
			growthFactor: 1.0,
			iterations:   10,
			expected:     1_024 * timeUnit,
		},
		{
			initial:      1 * timeUnit,
			growthFactor: 2.0,
			iterations:   5,
			expected:     243 * timeUnit,
		},
	}

	for caseIndex, c := range cases {
		t.Run(strconv.Itoa(caseIndex), func(t *testing.T) {
			g := gomega.NewWithT(t)
			actual := exponentialSleepTime(c.initial, c.growthFactor, c.iterations)
			g.Expect(actual).Should(gomega.Equal(c.expected))
		})
	}
}
