package ghost_test

import (
	"testing"

	"github.com/rliebz/ghost"
)

func TestBeNil(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		g := ghost.New(t)

		var v interface{}

		result := ghost.BeNil(v)()
		g.Should(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal("v is nil", result.Message))

		result = ghost.BeNil(nil)()
		g.Should(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal("nil is nil", result.Message))
	})

	t.Run("non-nil", func(t *testing.T) {
		g := ghost.New(t)

		var v int

		result := ghost.BeNil(v)()
		g.ShouldNot(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal("v is 0, not nil", result.Message))

		result = ghost.BeNil(-1 + 1)()
		g.ShouldNot(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal("-1 + 1 is 0, not nil", result.Message))
	})
}
