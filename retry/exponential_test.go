package retry

import (
	"testing"
	"time"
)

func TestExponential_Retry(t *testing.T) {
	cases := map[string]struct {
		config *Exponential
		retryOccurs
	}{
		"succeeds the first time it returns quickly": {
			config: &Exponential{
				InitialWaitBetweenAttempts: 1 * time.Second,
				GrowthFactor:               0.2,
			},
			retryOccurs: retryOccurs{
				errs:                  []error{Success},
				expectedDurationLower: time.Duration(0),
				expectedDurationUpper: 500 * time.Millisecond,
			},
		},
		"backs off and succeeds at 5 times": {
			config: &Exponential{
				InitialWaitBetweenAttempts: 10 * time.Millisecond,
				GrowthFactor:               1.0,
			},
			retryOccurs: retryOccurs{
				errs: []error{errAgain, errAgain, errAgain, errAgain, Success},
				// 10ms + 20ms + 40ms + 80ms = 150ms
				expectedDurationLower: 145 * time.Millisecond,
				expectedDurationUpper: 155 * time.Millisecond,
			},
		},
		"backs off and errors at 5 times": {
			config: &Exponential{
				InitialWaitBetweenAttempts: 10 * time.Millisecond,
				GrowthFactor:               1.0,
			},
			retryOccurs: retryOccurs{
				errs:        []error{errAgain, errAgain, errAgain, errAgain},
				expectedErr: errOutOfErrs,
				// 10ms + 20ms + 40ms + 80ms = 150ms
				expectedDurationLower: 145 * time.Millisecond,
				expectedDurationUpper: 155 * time.Millisecond,
			},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			c.Assert(t, c.config)
		})
	}
}
