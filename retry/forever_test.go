package retry

import (
	"github.com/wojnosystems/go-retry/retryStop"
	"testing"
	"time"
)

func TestForever_Retry(t *testing.T) {
	cases := map[string]struct {
		config *Forever
		retryOccurs
	}{
		"succeeds the first time it returns quickly": {
			config: &Forever{
				WaitBetweenAttempts: 1 * time.Second,
			},
			retryOccurs: retryOccurs{
				errs:                  []error{retryStop.Success},
				expectedDurationLower: time.Duration(0),
				expectedDurationUpper: 500 * time.Millisecond,
			},
		},
		"after 5 retries it succeeds": {
			config: &Forever{
				WaitBetweenAttempts: 10 * time.Millisecond,
			},
			retryOccurs: retryOccurs{
				errs: []error{errAgain, errAgain, errAgain, errAgain, retryStop.Success},
				// 10 + 10 + 10 + 10 = 40ms
				expectedDurationLower: 35 * time.Millisecond,
				expectedDurationUpper: 45 * time.Millisecond,
			},
		},
		"un-retryable error after 5 retries it errors": {
			config: &Forever{
				WaitBetweenAttempts: 10 * time.Millisecond,
			},
			retryOccurs: retryOccurs{
				errs:        []error{errAgain, errAgain, errAgain, errAgain, errAgain},
				expectedErr: errOutOfErrs,
				// 10 + 10 + 10 + 10 = 40ms
				expectedDurationLower: 45 * time.Millisecond,
				expectedDurationUpper: 55 * time.Millisecond,
			},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			c.Assert(t, c.config)
		})
	}
}
