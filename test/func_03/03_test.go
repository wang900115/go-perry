package func03_test

import (
	func03 "test/func_03"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDivide(t *testing.T) {
	result, err := func03.Divide(10, 2)
	require.NoError(t, err)
	assert.Equal(t, 5, result)

	_, err = func03.Divide(10, 0)
	require.Error(t, err)
}
