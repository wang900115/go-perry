package func04_test

import (
	func04 "test/func_04"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGreet(t *testing.T) {
	t.Run("with name", func(t *testing.T) {
		assert.Equal(t, "Hello, Alice", func04.Greet("Alice"))
	})

	t.Run("empty name", func(t *testing.T) {
		assert.Equal(t, "Hello, Guest", func04.Greet(""))
	})
}
