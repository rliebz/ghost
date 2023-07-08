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
	g.ShouldNot(be.False(true))

	g.Should(be.Equal(1+1, 2))
	g.Should(be.DeepEqual([]string{"a", "b"}, []string{"a", "b"}))

	g.Should(be.InSlice(2, []int{1, 2, 3}))
	g.Should(be.InString("foo", "foobar"))

	g.Should(be.Panic(func() { panic("oh no") }))
	g.ShouldNot(be.Panic(func() {}))

	var err error
	g.NoError(err)
	g.Must(be.Nil(err))
	g.MustNot(be.Error(err))

	err = errors.New("test error: oh no")
	g.Should(be.Error(err))
	g.Should(be.ErrorEqual("test error: oh no", err))
	g.Should(be.ErrorContaining("oh no", err))

	g.Should(be.JSONEqual(`{"b": 1, "a": 0}`, `{"a": 0, "b": 1}`))
	g.ShouldNot(be.JSONEqual(`{"a":1}`, `{"a":2}`))
}
