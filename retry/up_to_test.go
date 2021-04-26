package retry

import (
	"testing"
	"time"
)

func TestUpTo_Retry(t *testing.T) {
	cases := map[string]struct {
		config *UpTo
		retryOccurs
	}{
		"succeeds the first time it returns quickly": {
			config: &UpTo{
				WaitBetweenAttempts: 1 * time.Second,
				MaxAttempts:         5,
			},
			retryOccurs: retryOccurs{
				errs:                  []error{Success},
				expectedDurationLower: time.Duration(0),
				expectedDurationUpper: 500 * time.Millisecond,
			},
		},
		"after 5 retries it succeeds": {
			config: &UpTo{
				WaitBetweenAttempts: 10 * time.Millisecond,
				MaxAttempts:         6,
			},
			retryOccurs: retryOccurs{
				errs: []error{errAgain, errAgain, errAgain, errAgain, Success},
				// 10 + 10 + 10 + 10 = 40ms
				expectedDurationLower: 35 * time.Millisecond,
				expectedDurationUpper: 45 * time.Millisecond,
			},
		},
		"after 5 times it returns the retry error": {
			config: &UpTo{
				WaitBetweenAttempts: 10 * time.Millisecond,
				MaxAttempts:         5,
			},
			retryOccurs: retryOccurs{
				errs:        []error{errAgain, errAgain, errAgain, errAgain, errAgain},
				expectedErr: errAgain.wrapped,
				// 10 + 10 + 10 + 10 = 40ms
				expectedDurationLower: 35 * time.Millisecond,
				expectedDurationUpper: 45 * time.Millisecond,
			},
		},
		"after un-retryable error it returns the error": {
			config: &UpTo{
				WaitBetweenAttempts: 10 * time.Millisecond,
				MaxAttempts:         5,
			},
			retryOccurs: retryOccurs{
				errs:        []error{errAgain, errAgain, errFake},
				expectedErr: errFake,
				// 10 + 10 = 20ms
				expectedDurationLower: 15 * time.Millisecond,
				expectedDurationUpper: 25 * time.Millisecond,
			},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			c.Assert(t, c.config)
		})
	}
}
