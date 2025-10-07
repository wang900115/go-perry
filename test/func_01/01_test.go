package func_01_test

import (
	"test/func_01"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	assert.Equal(t, 5, func_01.Add(2, 3))
}
