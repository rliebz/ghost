package ghost_test

import (
	"errors"
	"testing"

	"github.com/rliebz/ghost"
	"github.com/rliebz/ghost/be"
)

func TestExample(t *testing.T) {
	g := ghost.New(t)

	g.Should(be.True(true))
	g.ShouldNot(be.True(false))

	g.Should(be.Equal(1+1, 2))
	g.Should(be.DeepEqual([]string{"a", "b"}, []string{"a", "b"}))
	g.Should(be.Containing([]int{1, 2, 3}, 2))
	g.Should(be.ContainingString("foobar", "foo"))

	g.Should(be.Panic(func() { panic("oh no") }))
	g.ShouldNot(be.Panic(func() {}))

	var err error
	g.NoError(err)
	g.Must(be.Nil(err))
	g.MustNot(be.Error(err))

	err = errors.New("oh my god")
	g.Should(be.ErrorContaining("my god", err))
	g.ShouldNot(be.ErrorContaining("steve", err))

	g.Should(be.JSONEqual(`{"b": 1, "a": 0}`, `{"a": 0, "b": 1}`))
	g.ShouldNot(be.JSONEqual(`{"a":1}`, `{"a":2}`))
}
