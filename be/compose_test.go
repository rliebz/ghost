package be_test

import (
	"testing"

	"github.com/rliebz/ghost"
	"github.com/rliebz/ghost/be"
)

func TestEventually(t *testing.T) {
	// TODO: Write tests
	_ = t
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
