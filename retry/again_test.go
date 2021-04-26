package retry

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAgain_Err(t *testing.T) {
	err := Again(errFake)
	a := err.(*again)
	require.Error(t, err)
	assert.EqualError(t, a.err(), err.Error())
}
