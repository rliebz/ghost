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
			return be.Equal(count, 3)
		}, 100*time.Millisecond, 5*time.Millisecond)

		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `count == 3`))
	})

	t.Run("not ok", func(t *testing.T) {
		count := 0
		result := be.Eventually(func() ghost.Result {
			count++
			return be.Equal(count, -1)
		}, 10*time.Millisecond, 5*time.Millisecond)

		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `count != -1
got:  1
want: -1
`))
	})

	t.Run("timeout", func(t *testing.T) {
		result := be.Eventually(func() ghost.Result {
			time.Sleep(100 * time.Millisecond)
			return be.True(true)
		}, 10*time.Millisecond, 100*time.Millisecond)

		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `func() ghost.Result {
	time.Sleep(100 * time.Millisecond)
	return be.True(true)
} did not return value within 10ms timeout`))
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
	g.Should(be.Equal(negated.Message, message))

	doubleNegated := be.Not(negated)
	g.Should(be.Equal(doubleNegated, result))
}
