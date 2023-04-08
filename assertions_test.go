package ghost_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/rliebz/ghost"
)

func TestBeNil(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		g := ghost.New(t)

		var v interface{}

		result := ghost.BeNil(v)()
		g.Should(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal("v is nil", result.Message))

		result = ghost.BeNil(nil)()
		g.Should(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal("nil is nil", result.Message))
	})

	t.Run("non-nil", func(t *testing.T) {
		g := ghost.New(t)

		var v int

		result := ghost.BeNil(v)()
		g.ShouldNot(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal("v is 0, not nil", result.Message))

		result = ghost.BeNil(-1 + 1)()
		g.ShouldNot(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal("-1 + 1 is 0, not nil", result.Message))
	})
}

func TestBeTrue(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		g := ghost.New(t)

		v := true
		result := ghost.BeTrue(v)()
		g.Should(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal("v is true", result.Message))

		result = ghost.BeTrue(true)()
		g.Should(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal("true is true", result.Message))
	})

	t.Run("false", func(t *testing.T) {
		g := ghost.New(t)

		v := false
		result := ghost.BeTrue(v)()
		g.ShouldNot(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal("v is false", result.Message))

		result = ghost.BeTrue(false)()
		g.ShouldNot(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal("false is false", result.Message))
	})
}

func TestBeZero(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		g := ghost.New(t)

		var v int
		result := ghost.BeZero(v)()
		g.Should(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal("v is the zero value", result.Message))

		result = ghost.BeZero(0)()
		g.Should(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal("0 is the zero value", result.Message))
	})

	t.Run("non-zero", func(t *testing.T) {
		g := ghost.New(t)

		v := 1
		result := ghost.BeZero(v)()
		g.ShouldNot(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal("v is non-zero\nvalue: 1", result.Message))

		result = ghost.BeZero(1)()
		g.ShouldNot(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal("1 is non-zero", result.Message))
	})
}

func TestContain(t *testing.T) {
	t.Run("contains <= 3", func(t *testing.T) {
		g := ghost.New(t)

		slice := []int{1, 2, 3}
		elem := 2

		result := ghost.Contain(slice, elem)()
		g.Should(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal(`slice contains elem
slice:   [1 2 3]
element: 2
`, result.Message))

		result = ghost.Contain([]int{1, 2, 3}, 2)()
		g.Should(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal(`[]int{1, 2, 3} contains 2
slice:   [1 2 3]
element: 2
`, result.Message))
	})

	t.Run("contains > 3", func(t *testing.T) {
		g := ghost.New(t)

		slice := []int{1, 2, 3, 4}
		elem := 2

		result := ghost.Contain(slice, elem)()
		g.Should(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal(`slice contains elem
slice:   [
	1
>	2
	3
	4
]
element: 2
`, result.Message))

		result = ghost.Contain([]int{1, 2, 3, 4}, 2)()
		g.Should(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal(`[]int{1, 2, 3, 4} contains 2
slice:   [
	1
>	2
	3
	4
]
element: 2
`, result.Message))
	})

	t.Run("does not contain <= 3", func(t *testing.T) {
		g := ghost.New(t)

		slice := []int{1, 2, 3}
		elem := 5

		result := ghost.Contain(slice, elem)()
		g.ShouldNot(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal(`slice does not contain elem
slice:   [1 2 3]
element: 5
`, result.Message))

		result = ghost.Contain([]int{1, 2, 3}, 5)()
		g.ShouldNot(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal(`[]int{1, 2, 3} does not contain 5
slice:   [1 2 3]
element: 5
`, result.Message))
	})

	t.Run("does not contain > 3", func(t *testing.T) {
		g := ghost.New(t)

		slice := []int{1, 2, 3, 4}
		elem := 5

		result := ghost.Contain(slice, elem)()
		g.ShouldNot(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal(`slice does not contain elem
slice:   [
	1
	2
	3
	4
]
element: 5
`, result.Message))

		result = ghost.Contain([]int{1, 2, 3, 4}, 5)()
		g.ShouldNot(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal(`[]int{1, 2, 3, 4} does not contain 5
slice:   [
	1
	2
	3
	4
]
element: 5
`, result.Message))
	})
}

func TestContainString(t *testing.T) {
	t.Run("contains", func(t *testing.T) {
		g := ghost.New(t)

		outer := "foobar"
		inner := "oob"

		result := ghost.ContainString(outer, inner)()
		g.Should(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal(`outer contains inner
str:    "foobar"
substr: "oob"
`, result.Message))

		result = ghost.ContainString("foobar", "oob")()
		g.Should(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal(`"foobar" contains "oob"
str:    "foobar"
substr: "oob"
`, result.Message))
	})

	t.Run("does not contain", func(t *testing.T) {
		g := ghost.New(t)

		outer := "foobar"
		inner := "boo"

		result := ghost.ContainString(outer, inner)()
		g.ShouldNot(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal(`outer does not contain inner
str:    "foobar"
substr: "boo"
`, result.Message))

		result = ghost.ContainString("foobar", "boo")()
		g.ShouldNot(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal(`"foobar" does not contain "boo"
str:    "foobar"
substr: "boo"
`, result.Message))
	})

	t.Run("multiline", func(t *testing.T) {
		g := ghost.New(t)

		outer := `one
two
three
`

		result := ghost.ContainString(outer, "two")()
		g.Should(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal(`outer contains "two"
str:    `+`
"""
one
two
three

"""
substr: "two"
`, result.Message))
	})
}

func TestDeepEqual(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		g := ghost.New(t)

		type T struct {
			A string
			B []int
		}

		want := T{"foo", []int{1, 2}}
		got := T{"foo", []int{1, 2}}

		result := ghost.DeepEqual(want, got)()
		g.Should(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal(`want == got
value: {foo [1 2]}`, result.Message))

		result = ghost.DeepEqual(T{"foo", []int{1}}, T{"foo", []int{1}})()
		g.Should(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal(`T{"foo", []int{1}} == T{"foo", []int{1}}
value: {foo [1]}`, result.Message))
	})

	t.Run("unequal", func(t *testing.T) {
		g := ghost.New(t)

		type T struct {
			A string
			B int
		}

		want := T{"foo", 1}
		got := T{"bar", 0}

		result := ghost.DeepEqual(want, got)()
		g.ShouldNot(ghost.BeTrue(result.Success))

		// Keep the diff small, because we don't want to test cmp.Diff
		wantText := `want != got
diff (-want +got):
  ghost_test.T{
- 	A: "foo",
+ 	A: "bar",
- 	B: 1,
+ 	B: 0,
  }

`
		result.Message = strings.ReplaceAll(result.Message, "\u00a0", " ")
		g.Should(ghost.Equal(wantText, result.Message))

		result = ghost.DeepEqual(T{"foo", 1}, T{"bar", 0})()
		g.ShouldNot(ghost.BeTrue(result.Success))

		wantText = `T{"foo", 1} != T{"bar", 0}
diff (-want +got):
  ghost_test.T{
- 	A: "foo",
+ 	A: "bar",
- 	B: 1,
+ 	B: 0,
  }

`
		result.Message = strings.ReplaceAll(result.Message, "\u00a0", " ")
		g.Should(ghost.Equal(wantText, result.Message))
	})
}

func TestEqual(t *testing.T) {
	// TODO
}

func TestError(t *testing.T) {
	// TODO
}

func TestErrorContaining(t *testing.T) {
	t.Run("contains", func(t *testing.T) {
		g := ghost.New(t)

		err := errors.New("foobar")
		msg := "oob"

		result := ghost.ErrorContaining(err, msg)()
		g.Should(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal(`error err contains message "oob": foobar`, result.Message))

		result = ghost.ErrorContaining(errors.New("foobar"), "oob")()
		g.Should(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal(`error errors.New("foobar") contains message "oob": foobar`, result.Message))
	})

	t.Run("does not contain", func(t *testing.T) {
		g := ghost.New(t)

		err := errors.New("foobar")
		msg := "boo"

		result := ghost.ErrorContaining(err, msg)()
		g.ShouldNot(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal(`error err does not contain message "boo": foobar`, result.Message))

		result = ghost.ErrorContaining(errors.New("foobar"), "boo")()
		g.ShouldNot(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal(`error errors.New("foobar") does not contain message "boo": foobar`, result.Message))
	})

	t.Run("nil", func(t *testing.T) {
		g := ghost.New(t)

		var err error
		msg := "boo"

		result := ghost.ErrorContaining(err, msg)()
		g.ShouldNot(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal(`error err is nil; missing message msg: boo`, result.Message))

		result = ghost.ErrorContaining(nil, "boo")()
		g.ShouldNot(ghost.BeTrue(result.Success))
		g.Should(ghost.Equal(`error nil is nil; missing message: boo`, result.Message))
	})
}
