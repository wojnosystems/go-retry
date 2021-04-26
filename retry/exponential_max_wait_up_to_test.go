package retry

import (
	"testing"
	"time"
)

func TestExponentialMaxWaitUpTo_Retry(t *testing.T) {
	cases := map[string]struct {
		config *ExponentialMaxWaitUpTo
		retryOccurs
	}{
		"succeeds the first time it returns quickly": {
			config: &ExponentialMaxWaitUpTo{
				InitialWaitBetweenAttempts: 1 * time.Second,
				GrowthFactor:               0.2,
				MaxAttempts:                10,
				MaxWaitBetweenAttempts:     2 * time.Second,
			},
			retryOccurs: retryOccurs{
				errs:                  []error{Success},
				expectedDurationLower: time.Duration(0),
				expectedDurationUpper: 500 * time.Millisecond,
			},
		},
		"backs off and succeeds at 5 times": {
			config: &ExponentialMaxWaitUpTo{
				InitialWaitBetweenAttempts: 10 * time.Millisecond,
				GrowthFactor:               1.0,
				MaxAttempts:                10,
				MaxWaitBetweenAttempts:     15 * time.Millisecond,
			},
			retryOccurs: retryOccurs{
				errs: []error{errAgain, errAgain, errAgain, errAgain, Success},
				// 10ms + 15ms + 15ms + 15ms = 55ms
				expectedDurationLower: 50 * time.Millisecond,
				expectedDurationUpper: 60 * time.Millisecond,
			},
		},
		"backs off and runs out of retries": {
			config: &ExponentialMaxWaitUpTo{
				InitialWaitBetweenAttempts: 10 * time.Millisecond,
				GrowthFactor:               1.0,
				MaxAttempts:                5,
				MaxWaitBetweenAttempts:     15 * time.Millisecond,
			},
			retryOccurs: retryOccurs{
				errs:        []error{errAgain, errAgain, errAgain, errAgain, errAgain},
				expectedErr: errAgain.err(),
				// 10ms + 15ms + 15ms + 15ms = 55ms
				expectedDurationLower: 50 * time.Millisecond,
				expectedDurationUpper: 60 * time.Millisecond,
			},
		},
		"backs off and errors at 5": {
			config: &ExponentialMaxWaitUpTo{
				InitialWaitBetweenAttempts: 10 * time.Millisecond,
				GrowthFactor:               1.0,
				MaxAttempts:                6,
				MaxWaitBetweenAttempts:     15 * time.Millisecond,
			},
			retryOccurs: retryOccurs{
				errs:        []error{errAgain, errAgain, errAgain, errAgain, errAgain},
				expectedErr: errOutOfErrs,
				// 10ms + 15ms + 15ms + 15ms + 15ms = 70ms
				expectedDurationLower: 65 * time.Millisecond,
				expectedDurationUpper: 75 * time.Millisecond,
			},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			c.Assert(t, c.config)
		})
	}
}
