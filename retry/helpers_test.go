package retry

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/wojnosystems/go-retry/examples/common"
	"github.com/wojnosystems/go-retry/retryAgain"
	"testing"
	"time"
)

var errFake = errors.New("fake")
var errAgain = retryAgain.Error(errors.New("again"))
var errOutOfErrs = errors.New("ran out of errors")

type retryOccurs struct {
	errs                  []error
	expectedDurationLower time.Duration
	expectedDurationUpper time.Duration
	expectedErr           error
}

func (o *retryOccurs) Assert(t *testing.T, retrier Retrier) {
	tries := 0
	duration := common.TimeThis(func() {
		err := retrier.Retry(func() (err error) {
			defer func() {
				tries++
			}()
			if len(o.errs)-1 < tries {
				return errOutOfErrs
			}
			return o.errs[tries]
		})
		if o.expectedErr != nil {
			assert.EqualError(t, err, o.expectedErr.Error())
		} else {
			assert.NoError(t, err)
		}
	})
	assert.Greater(t, duration, o.expectedDurationLower)
	assert.Less(t, duration, o.expectedDurationUpper)
}
