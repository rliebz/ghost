package ghost_test

import (
	"errors"
	"testing"

	"github.com/rliebz/ghost"
)

func TestExample(t *testing.T) {
	g := ghost.New(t)

	g.Should(ghost.BeTrue(true))
	g.ShouldNot(ghost.BeTrue(false))

	g.Should(ghost.Equal(1+1, 2))
	g.Should(ghost.DeepEqual([]string{"a", "b"}, []string{"a", "b"}))
	g.Should(ghost.Contain([]int{1, 2, 3}, 2))
	g.Should(ghost.ContainString("foobar", "foo"))

	g.Should(ghost.Panic(func() { panic("oh no") }))
	g.ShouldNot(ghost.Panic(func() {}))

	var err error
	g.NoError(err)
	g.Must(ghost.BeNil(err))
	g.MustNot(ghost.Error(err))

	err = errors.New("oh my god")
	g.Should(ghost.ErrorContaining("my god", err))
	g.ShouldNot(ghost.ErrorContaining("steve", err))

	g.Should(ghost.JSONEqual(`{"b": 1, "a": 0}`, `{"a": 0, "b": 1}`))
	g.ShouldNot(ghost.JSONEqual(`{"a":1}`, `{"a":2}`))
}
