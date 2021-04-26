package retry

import (
	"testing"
	"time"
)

func TestExponentialUpTo_Retry(t *testing.T) {
	cases := map[string]struct {
		config *ExponentialUpTo
		retryOccurs
	}{
		"succeeds the first time it returns quickly": {
			config: &ExponentialUpTo{
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
			config: &ExponentialUpTo{
				InitialWaitBetweenAttempts: 10 * time.Millisecond,
				GrowthFactor:               0.5,
				MaxAttempts:                10,
			},
			retryOccurs: retryOccurs{
				errs: []error{errAgain, errAgain, errAgain, errAgain, Success},
				// 10ms + 15ms + 22.5ms + 33.75ms = 81.25ms
				expectedDurationLower: 76 * time.Millisecond,
				expectedDurationUpper: 87 * time.Millisecond,
			},
		},
		"backs off and runs out of retries": {
			config: &ExponentialUpTo{
				InitialWaitBetweenAttempts: 10 * time.Millisecond,
				GrowthFactor:               0.5,
				MaxAttempts:                5,
			},
			retryOccurs: retryOccurs{
				errs:        []error{errAgain, errAgain, errAgain, errAgain, errAgain},
				expectedErr: errAgain.err(),
				// 10ms + 15ms + 22.5ms + 33.75ms = 81.25ms
				expectedDurationLower: 76 * time.Millisecond,
				expectedDurationUpper: 87 * time.Millisecond,
			},
		},
		"backs off and fails at 5 times": {
			config: &ExponentialUpTo{
				InitialWaitBetweenAttempts: 10 * time.Millisecond,
				GrowthFactor:               0.5,
				MaxAttempts:                5,
			},
			retryOccurs: retryOccurs{
				errs:        []error{errAgain, errAgain, errAgain, errAgain},
				expectedErr: errOutOfErrs,
				// 10ms + 15ms + 22.5ms + 33.75ms = 81.25ms
				expectedDurationLower: 76 * time.Millisecond,
				expectedDurationUpper: 87 * time.Millisecond,
			},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			c.Assert(t, c.config)
		})
	}
}
