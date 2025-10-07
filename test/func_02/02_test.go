package func_02_test

import (
	"test/func_02"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEven(t *testing.T) {
	tests := []struct {
		input int
		want  bool
	}{
		{2, true},
		{3, false},
		{0, true},
		{-4, true},
		{-1, false},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, func_02.IsEven(tt.input))
	}
}
