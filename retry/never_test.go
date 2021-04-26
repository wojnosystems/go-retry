package retry

import (
	"github.com/wojnosystems/go-retry/retryStop"
	"testing"
	"time"
)

func TestNever_Retry(t *testing.T) {
	cases := map[string]struct {
		config *UpTo
		retryOccurs
	}{
		"succeeds the first time it returns quickly": {
			config: Never,
			retryOccurs: retryOccurs{
				errs:                  []error{retryStop.Success},
				expectedDurationLower: time.Duration(0),
				expectedDurationUpper: 500 * time.Millisecond,
			},
		},
		"fails once": {
			config: Never,
			retryOccurs: retryOccurs{
				errs:                  []error{errAgain, errAgain},
				expectedErr:           errAgain,
				expectedDurationLower: 0,
				expectedDurationUpper: 10 * time.Millisecond,
			},
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			c.Assert(t, c.config)
		})
	}
}
