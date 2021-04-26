package core

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/wojnosystems/go-retry/retryAgain"
	"github.com/wojnosystems/go-retry/retryStop"
	"testing"
)

var errAgain = retryAgain.Error(errors.New("retry-again"))
var errRetriesExceeded = errors.New("no more retries")

func TestLoopForever(t *testing.T) {
	cases := map[string]struct {
		returns       []error
		expectedCount uint64
		expectedErr   error
	}{
		"one success it returns": {
			returns: []error{retryStop.Success},
		},
		"retries occur": {
			returns:       []error{errAgain, errAgain, retryStop.Success},
			expectedCount: 2,
		},
		"retries exceeded": {
			returns:       []error{errAgain, errAgain, errRetriesExceeded},
			expectedCount: 2,
			expectedErr:   errRetriesExceeded,
		},
	}

	for caseName, c := range cases {
		t.Run(caseName, func(t *testing.T) {
			count := uint64(0)
			actual := LoopForever(func() (err error) {
				if count > uint64(len(c.returns)) {
					return
				}
				return c.returns[count]
			}, func(iterator uint64) {
				count = iterator + 1
			})
			if c.expectedErr != nil {
				assert.EqualError(t, actual, c.expectedErr.Error())
			} else {
				assert.NoError(t, actual)
			}
		})
	}
}
