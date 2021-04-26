package retry

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/wojnosystems/go-retry/retryAgain"
	"testing"
)

func TestSkip_Retry(t *testing.T) {
	called := false
	skipErr := (&Skip{}).Retry(context.Background(), func() (err error) {
		called = true
		return retryAgain.Error(errors.New("fake"))
	})
	assert.NoError(t, skipErr)
	assert.False(t, called)
}
