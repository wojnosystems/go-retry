package core

import (
	"github.com/stretchr/testify/assert"
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
			actual := MinDuration(c.a, c.b)
			assert.Equal(t, c.expected, actual)
		})
	}
}
