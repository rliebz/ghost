package be_test

import (
	"testing"

	"github.com/rliebz/ghost"
	"github.com/rliebz/ghost/be"
)

func TestGreater(t *testing.T) {
	t.Run("int less", func(t *testing.T) {
		g := ghost.New(t)

		a := 3
		b := 4

		result := be.Greater(a, b)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `a (3) is not greater than b (4)`))

		result = be.Greater(3, 4)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `3 is not greater than 4`))
	})

	t.Run("int greater", func(t *testing.T) {
		g := ghost.New(t)

		a := 4
		b := 3

		result := be.Greater(a, b)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `a (4) is greater than b (3)`))

		result = be.Greater(4, 3)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `4 is greater than 3`))
	})

	t.Run("int equal", func(t *testing.T) {
		g := ghost.New(t)

		a := 3
		b := 3

		result := be.Greater(a, b)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `a (3) is not greater than b (3)`))

		result = be.Greater(3, 3)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `3 is not greater than 3`))
	})

	t.Run("float less", func(t *testing.T) {
		g := ghost.New(t)

		a := 3.1
		b := 4.1

		result := be.Greater(a, b)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `a (3.1) is not greater than b (4.1)`))

		result = be.Greater(3.1, 4.1)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `3.1 is not greater than 4.1`))
	})

	t.Run("float greater", func(t *testing.T) {
		g := ghost.New(t)

		a := 4.1
		b := 3.1

		result := be.Greater(a, b)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `a (4.1) is greater than b (3.1)`))

		result = be.Greater(4.1, 3.1)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `4.1 is greater than 3.1`))
	})

	t.Run("float equal", func(t *testing.T) {
		g := ghost.New(t)

		a := 3.1
		b := 3.1

		result := be.Greater(a, b)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `a (3.1) is not greater than b (3.1)`))

		result = be.Greater(3.1, 3.1)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `3.1 is not greater than 3.1`))
	})

	t.Run("string less", func(t *testing.T) {
		g := ghost.New(t)

		a := "bar"
		b := "foo"

		result := be.Greater(a, b)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `a ("bar") is not greater than b ("foo")`))

		result = be.Greater("bar", "foo")
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `"bar" is not greater than "foo"`))
	})

	t.Run("string greater", func(t *testing.T) {
		g := ghost.New(t)

		a := "foo"
		b := "bar"

		result := be.Greater(a, b)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `a ("foo") is greater than b ("bar")`))

		result = be.Greater("foo", "bar")
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `"foo" is greater than "bar"`))
	})

	t.Run("string equal", func(t *testing.T) {
		g := ghost.New(t)

		a := "foo"
		b := "foo"

		result := be.Greater(a, b)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `a ("foo") is not greater than b ("foo")`))

		result = be.Greater("foo", "foo")
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `"foo" is not greater than "foo"`))
	})
}

func TestLess(t *testing.T) {
	t.Run("int less", func(t *testing.T) {
		g := ghost.New(t)

		a := 3
		b := 4

		result := be.Less(a, b)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `a (3) is less than b (4)`))

		result = be.Less(3, 4)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `3 is less than 4`))
	})

	t.Run("int greater", func(t *testing.T) {
		g := ghost.New(t)

		a := 4
		b := 3

		result := be.Less(a, b)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `a (4) is not less than b (3)`))

		result = be.Less(4, 3)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `4 is not less than 3`))
	})

	t.Run("int equal", func(t *testing.T) {
		g := ghost.New(t)

		a := 3
		b := 3

		result := be.Less(a, b)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `a (3) is not less than b (3)`))

		result = be.Less(3, 3)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `3 is not less than 3`))
	})

	t.Run("float less", func(t *testing.T) {
		g := ghost.New(t)

		a := 3.1
		b := 4.1

		result := be.Less(a, b)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `a (3.1) is less than b (4.1)`))

		result = be.Less(3.1, 4.1)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `3.1 is less than 4.1`))
	})

	t.Run("float greater", func(t *testing.T) {
		g := ghost.New(t)

		a := 4.1
		b := 3.1

		result := be.Less(a, b)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `a (4.1) is not less than b (3.1)`))

		result = be.Less(4.1, 3.1)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `4.1 is not less than 3.1`))
	})

	t.Run("float equal", func(t *testing.T) {
		g := ghost.New(t)

		a := 3.1
		b := 3.1

		result := be.Less(a, b)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `a (3.1) is not less than b (3.1)`))

		result = be.Less(3.1, 3.1)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `3.1 is not less than 3.1`))
	})

	t.Run("string less", func(t *testing.T) {
		g := ghost.New(t)

		a := "bar"
		b := "foo"

		result := be.Less(a, b)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `a ("bar") is less than b ("foo")`))

		result = be.Less("bar", "foo")
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `"bar" is less than "foo"`))
	})

	t.Run("string greater", func(t *testing.T) {
		g := ghost.New(t)

		a := "foo"
		b := "bar"

		result := be.Less(a, b)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `a ("foo") is not less than b ("bar")`))

		result = be.Less("foo", "bar")
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `"foo" is not less than "bar"`))
	})

	t.Run("string equal", func(t *testing.T) {
		g := ghost.New(t)

		a := "foo"
		b := "foo"

		result := be.Less(a, b)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `a ("foo") is not less than b ("foo")`))

		result = be.Less("foo", "foo")
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `"foo" is not less than "foo"`))
	})
}
