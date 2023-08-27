package be_test

import (
	"testing"
	"time"

	"github.com/rliebz/ghost"
	"github.com/rliebz/ghost/be"
)

func TestEventually(t *testing.T) {
	g := ghost.New(t)

	t.Run("ok", func(t *testing.T) {
		count := 0
		result := be.Eventually(func() ghost.Result {
			count++
			return be.Equal(3, count)
		}, 100*time.Millisecond, 5*time.Millisecond)

		g.Should(be.True(result.Ok))
		g.Should(be.Equal(`3 == count`, result.Message))
	})

	t.Run("not ok", func(t *testing.T) {
		count := 0
		result := be.Eventually(func() ghost.Result {
			count++
			return be.Equal(-1, count)
		}, 10*time.Millisecond, 5*time.Millisecond)

		g.Should(be.False(result.Ok))
		g.Should(be.Equal(`-1 != count
want: -1
got:  1
`, result.Message))
	})

	t.Run("timeout", func(t *testing.T) {
		result := be.Eventually(func() ghost.Result {
			time.Sleep(100 * time.Millisecond)
			return be.True(true)
		}, 10*time.Millisecond, 100*time.Millisecond)

		g.Should(be.False(result.Ok))
		g.Should(be.Equal(`func() ghost.Result {
	time.Sleep(100 * time.Millisecond)
	return be.True(true)
} did not return value within 10ms timeout`, result.Message))
	})
}

func TestNot(t *testing.T) {
	g := ghost.New(t)

	message := "some message"

	result := ghost.Result{
		Ok:      true,
		Message: message,
	}

	negated := be.Not(result)
	g.Should(be.False(negated.Ok))
	g.Should(be.Equal(message, negated.Message))

	doubleNegated := be.Not(negated)
	g.Should(be.Equal(result, doubleNegated))
}
