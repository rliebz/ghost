package ghost_test

import (
	"errors"
	"testing"
	"time"

	"github.com/rliebz/ghost"
	"github.com/rliebz/ghost/be"
)

func TestExample(t *testing.T) {
	g := ghost.New(t)

	g.Should(be.True(true))
	g.Should(be.Not(be.False(true)))

	g.Should(be.Equal(1+1, 2))
	g.Should(be.DeepEqual([]string{"a", "b"}, []string{"a", "b"}))

	g.Should(be.Close(11.0, 6.8+4.2, 0.01))

	g.Should(be.MapLen(2, map[string]int{"a": 1, "b": 2}))

	g.Should(be.SliceContaining("b", []string{"a", "b", "c"}))
	g.Should(be.SliceLen(3, []string{"a", "b", "c"}))

	g.Should(be.StringContaining("foo", "foobar"))

	g.Should(be.Panic(func() { panic("oh no") }))
	g.Should(be.Not(be.Panic(func() {})))

	var err error
	g.NoError(err)
	g.Must(be.Nil(err))
	g.Must(be.Not(be.Error(err)))

	err = errors.New("test error: oh no")
	g.Should(be.Error(err))
	g.Should(be.ErrorEqual("test error: oh no", err))
	g.Should(be.ErrorContaining("oh no", err))

	g.Should(be.JSONEqual(`{"b": 1, "a": 0}`, `{"a": 0, "b": 1}`))
	g.Should(be.Not(be.JSONEqual(`{"a":1}`, `{"a":2}`)))

	count := 0
	g.Should(be.Eventually(func() ghost.Result {
		count++
		return be.Equal(3, count)
	}, 100*time.Millisecond, 10*time.Millisecond))
}
