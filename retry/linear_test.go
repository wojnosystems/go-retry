package retry

import (
	"github.com/wojnosystems/go-retry/retryStop"
	"testing"
	"time"
)

func TestLinear_Retry(t *testing.T) {
	cases := map[string]struct {
		config *Linear
		retryOccurs
	}{
		"succeeds the first time it returns quickly": {
			config: &Linear{
				InitialWaitBetweenAttempts: 1 * time.Second,
				GrowthFactor:               0.2,
			},
			retryOccurs: retryOccurs{
				errs:                  []error{retryStop.Success},
				expectedDurationLower: time.Duration(0),
				expectedDurationUpper: 500 * time.Millisecond,
			},
		},
		"backs off and succeeds at 5 times": {
			config: &Linear{
				InitialWaitBetweenAttempts: 10 * time.Millisecond,
				GrowthFactor:               0.2,
			},
			retryOccurs: retryOccurs{
				errs: []error{errAgain, errAgain, errAgain, errAgain, retryStop.Success},
				// 10ms + 12ms + 14.4ms + 17.28ms = 53.68ms
				expectedDurationLower: 48 * time.Millisecond,
				expectedDurationUpper: 58 * time.Millisecond,
			},
		},
		"backs off and errors at 5 times": {
			config: &Linear{
				InitialWaitBetweenAttempts: 10 * time.Millisecond,
				GrowthFactor:               0.2,
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
