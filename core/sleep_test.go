package core

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSleep(t *testing.T) {
	cases := map[string]struct {
		sleepDuration time.Duration
		makeContext   func(cb func(ctx context.Context))
		expectedLower time.Duration
		expectedUpper time.Duration
	}{
		"sleep is longer": {
			sleepDuration: 10 * time.Second,
			makeContext: func(cb func(ctx context.Context)) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
				defer cancel()
				cb(ctx)
			},
			expectedLower: 5 * time.Millisecond,
			expectedUpper: 15 * time.Millisecond,
		},
		"ctx is longer": {
			sleepDuration: 10 * time.Millisecond,
			makeContext: func(cb func(ctx context.Context)) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				cb(ctx)
			},
			expectedLower: 5 * time.Millisecond,
			expectedUpper: 15 * time.Millisecond,
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			start := time.Now()
			c.makeContext(func(ctx context.Context) {
				Sleep(ctx, c.sleepDuration)
			})
			duration := time.Now().Sub(start)
			assert.Greater(t, duration, c.expectedLower)
			assert.Less(t, duration, c.expectedUpper)
		})
	}
}
