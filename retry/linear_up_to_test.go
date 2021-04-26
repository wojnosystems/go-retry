package retry

import (
	"testing"
	"time"
)

func TestLinearUpTo_Retry(t *testing.T) {
	cases := map[string]struct {
		config *LinearUpTo
		retryOccurs
	}{
		"succeeds the first time it returns quickly": {
			config: &LinearUpTo{
				InitialWaitBetweenAttempts: 1 * time.Second,
				GrowthFactor:               0.2,
				MaxAttempts:                10,
			},
			retryOccurs: retryOccurs{
				errs:                  []error{Success},
				expectedDurationLower: time.Duration(0),
				expectedDurationUpper: 500 * time.Millisecond,
			},
		},
		"backs off and retries 5 times then succeeds": {
			config: &LinearUpTo{
				InitialWaitBetweenAttempts: 10 * time.Millisecond,
				GrowthFactor:               0.2,
				MaxAttempts:                10,
			},
			retryOccurs: retryOccurs{
				errs: []error{errAgain, errAgain, errAgain, errAgain, Success},
				// 10ms + 12ms + 14.4ms + 17.28ms = 53.68ms
				expectedDurationLower: 48 * time.Millisecond,
				expectedDurationUpper: 58 * time.Millisecond,
			},
		},
		"backs off and runs out of retries": {
			config: &LinearUpTo{
				InitialWaitBetweenAttempts: 10 * time.Millisecond,
				GrowthFactor:               0.2,
				MaxAttempts:                5,
			},
			retryOccurs: retryOccurs{
				errs:        []error{errAgain, errAgain, errAgain, errAgain, errAgain},
				expectedErr: errAgain.err(),
				// 10ms + 12ms + 14.4ms + 17.28ms = 53.68ms
				expectedDurationLower: 48 * time.Millisecond,
				expectedDurationUpper: 58 * time.Millisecond,
			},
		},
		"backs off and fails at 5 times": {
			config: &LinearUpTo{
				InitialWaitBetweenAttempts: 10 * time.Millisecond,
				GrowthFactor:               0.2,
				MaxAttempts:                5,
			},
			retryOccurs: retryOccurs{
				errs:        []error{errAgain, errAgain, errAgain, errAgain},
				expectedErr: errOutOfErrs,
				// 10ms + 12ms + 14.4ms + 17.28ms = 53.68ms
				expectedDurationLower: 48 * time.Millisecond,
				expectedDurationUpper: 58 * time.Millisecond,
			},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			c.Assert(t, c.config)
		})
	}
}
