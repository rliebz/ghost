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
		g.Should(be.Eventually(func() ghost.Result {
			count++
			return be.Equal(3, count)
		}, 100*time.Millisecond, 5*time.Millisecond))
	})

	t.Run("not ok", func(t *testing.T) {
		count := 0
		g.Should(be.Not(be.Eventually(func() ghost.Result {
			count++
			return be.Equal(-1, count)
		}, 10*time.Millisecond, 5*time.Millisecond)))
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
