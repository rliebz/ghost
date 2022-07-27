package ghost_test

import (
	"testing"

	"github.com/rliebz/ghost"
)

func TestExample(t *testing.T) {
	it := ghost.New(t)

	it.Should(ghost.BeTrue(true))
	it.ShouldNot(ghost.BeTrue(false))

	it.Should(ghost.Equal(1+1, 2))
	it.Should(ghost.DeepEqual([]string{"a", "b"}, []string{"a", "b"}))
	it.Should(ghost.Contain([]int{1, 2, 3}, 2))
	it.Should(ghost.ContainString("foobar", "foo"))

	it.Should(ghost.Panic(func() { panic("oh no") }))
	it.ShouldNot(ghost.Panic(func() {}))

	var err error
	it.MustNot(ghost.Err(err))
}
