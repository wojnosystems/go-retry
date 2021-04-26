package core

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/wojnosystems/go-retry/retryStop"
	"testing"
)

var errFake = errors.New("fake error")

func TestLoopUpTo(t *testing.T) {
	cases := map[string]struct {
		returns       []error
		expectedCount uint64
		expectedErr   error
		maxAttempts   uint64
	}{
		"tries once": {
			returns:       []error{errAgain, errAgain, retryStop.Success},
			expectedCount: 0,
			expectedErr:   errAgain,
		},
		"one success it returns": {
			returns:     []error{retryStop.Success},
			maxAttempts: 100,
		},
		"retries occur": {
			returns:       []error{errAgain, errAgain, retryStop.Success},
			expectedCount: 2,
			maxAttempts:   100,
		},
		"retries exceeded": {
			returns:       []error{errAgain, errAgain, errRetriesExceeded},
			expectedCount: 2,
			expectedErr:   errAgain,
			maxAttempts:   1,
		},
		"un-retryable error": {
			returns:       []error{errAgain, errAgain, errFake, errAgain},
			expectedCount: 2,
			expectedErr:   errFake,
			maxAttempts:   100,
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			count := uint64(0)
			actual := LoopUpTo(func() (err error) {
				if count > uint64(len(c.returns)) {
					return
				}
				return c.returns[count]
			}, func(iterator uint64) {
				count = iterator + 1
			}, c.maxAttempts)
			if c.expectedErr != nil {
				assert.EqualError(t, actual, c.expectedErr.Error())
			} else {
				assert.NoError(t, actual)
			}
		})
	}
}
