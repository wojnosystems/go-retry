package retryAgain

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var errFake = errors.New("fake")

func TestAgain_Err(t *testing.T) {
	err := Error(errFake)
	require.Error(t, err)
	assert.EqualError(t, errFake, err.Error())
}
