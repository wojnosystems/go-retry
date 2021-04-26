package retry

import (
	"github.com/wojnosystems/go-retry/retryStop"
	"testing"
	"time"
)

func TestLinearMaxWaitUpTo_Retry(t *testing.T) {
	cases := map[string]struct {
		config *LinearMaxWaitUpTo
		retryOccurs
	}{
		"succeeds the first time it returns quickly": {
			config: &LinearMaxWaitUpTo{
				InitialWaitBetweenAttempts: 1 * time.Second,
				GrowthFactor:               0.2,
				MaxAttempts:                10,
				MaxWaitBetweenAttempts:     2 * time.Second,
			},
			retryOccurs: retryOccurs{
				errs:                  []error{retryStop.Success},
				expectedDurationLower: time.Duration(0),
				expectedDurationUpper: 500 * time.Millisecond,
			},
		},
		"backs off and succeeds at 5 times": {
			config: &LinearMaxWaitUpTo{
				InitialWaitBetweenAttempts: 10 * time.Millisecond,
				GrowthFactor:               1.0,
				MaxAttempts:                10,
				MaxWaitBetweenAttempts:     15 * time.Millisecond,
			},
			retryOccurs: retryOccurs{
				errs: []error{errAgain, errAgain, errAgain, errAgain, retryStop.Success},
				// 10ms + 15ms + 15.0ms + 15.0ms = 55ms
				expectedDurationLower: 50 * time.Millisecond,
				expectedDurationUpper: 60 * time.Millisecond,
			},
		},
		"backs off and runs out of retries": {
			config: &LinearMaxWaitUpTo{
				InitialWaitBetweenAttempts: 10 * time.Millisecond,
				GrowthFactor:               0.2,
				MaxAttempts:                5,
				MaxWaitBetweenAttempts:     15 * time.Millisecond,
			},
			retryOccurs: retryOccurs{
				errs:        []error{errAgain, errAgain, errAgain, errAgain, errAgain},
				expectedErr: errAgain.Err(),
				// 10ms + 15ms + 15ms + 15ms = 55ms
				expectedDurationLower: 50 * time.Millisecond,
				expectedDurationUpper: 60 * time.Millisecond,
			},
		},
		"backs off and errors at 5": {
			config: &LinearMaxWaitUpTo{
				InitialWaitBetweenAttempts: 10 * time.Millisecond,
				GrowthFactor:               0.2,
				MaxAttempts:                6,
				MaxWaitBetweenAttempts:     15 * time.Millisecond,
			},
			retryOccurs: retryOccurs{
				errs:        []error{errAgain, errAgain, errAgain, errAgain, errAgain},
				expectedErr: errOutOfErrs,
				// 10ms + 15ms + 15ms + 15ms + 15 = 70ms
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
