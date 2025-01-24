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
	g.ShouldNot(be.False(true))

	g.Should(be.Equal(1+1, 2))
	g.Should(be.DeepEqual([]string{"a", "b"}, []string{"a", "b"}))

	g.Should(be.Close(6.8+4.2, 11.0, 0.01))

	g.Should(be.MapLen(map[string]int{"a": 1, "b": 2}, 2))

	g.Should(be.SliceContaining([]string{"a", "b", "c"}, "b"))
	g.Should(be.SliceLen([]string{"a", "b", "c"}, 3))

	g.Should(be.StringContaining("foobar", "foo"))
	g.Should(be.StringMatching("foobar", `^foo`))

	var err error
	g.NoError(err)
	g.Must(be.Nil(err))
	g.MustNot(be.Error(err))

	err = errors.New("test error: oh no")
	g.Should(be.Error(err))
	g.Should(be.ErrorEqual(err, "test error: oh no"))
	g.Should(be.ErrorContaining(err, "oh no"))

	g.Should(be.JSONEqual(`{"b": 1, "a": 0}`, `{"a": 0, "b": 1}`))
	g.ShouldNot(be.JSONEqual(`{"a":1}`, `{"a":2}`))
}

func ExampleAssignedAs() {
	t := new(testing.T) // from the test
	g := ghost.New(t)

	defer func() {
		var err error
		g.Must(be.AssignedAs(recover(), &err))
		g.Should(be.ErrorEqual(err, "oops"))
	}()

	panic(errors.New("oops"))
}

func ExampleEventually() {
	t := new(testing.T) // from the test
	g := ghost.New(t)

	count := 0
	g.Should(be.Eventually(func() ghost.Result {
		count++
		return be.Equal(count, 3)
	}, 100*time.Millisecond, 10*time.Millisecond))
}

func ExampleAny() {
	t := new(testing.T) // from the test
	g := ghost.New(t)

	g.Should(be.Any(
		be.Equal(0, 1),
		be.Equal(1, 1),
		be.Equal(2, 1),
	))
}

func ExampleAll() {
	t := new(testing.T) // from the test
	g := ghost.New(t)

	g.Should(be.All(
		be.Equal(1, 1),
		be.Equal(2, 2),
		be.Equal(3, 3),
	))
}
