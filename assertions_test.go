package ghost_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/rliebz/ghost"
)

func TestBeInDelta(t *testing.T) {
	t.Run("in delta", func(t *testing.T) {
		g := ghost.New(t)

		want := 32.5
		got := 32.0

		result := ghost.BeInDelta(want, got, 1)
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal(
			"delta 0.5 between want (32.5) and got (32) is within 1",
			result.Message,
		))

		result = ghost.BeInDelta(32.5, 32.0, 1.0)
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal(
			"delta 0.5 between 32.5 and 32.0 is within 1",
			result.Message,
		))

		result = ghost.BeInDelta(32.0, 32.5, 1.0)
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal(
			"delta 0.5 between 32.0 and 32.5 is within 1",
			result.Message,
		))
	})

	t.Run("not in delta", func(t *testing.T) {
		g := ghost.New(t)

		want := 32.5
		got := 32.0

		result := ghost.BeInDelta(want, got, 0.3)
		g.Should(ghost.BeFalse(result.Ok))
		g.Should(ghost.Equal(
			"delta 0.5 between want (32.5) and got (32) is not within 0.3",
			result.Message,
		))

		result = ghost.BeInDelta(32.5, 32.0, 0.3)
		g.Should(ghost.BeFalse(result.Ok))
		g.Should(ghost.Equal(
			"delta 0.5 between 32.5 and 32.0 is not within 0.3",
			result.Message,
		))
	})
}

func TestBeNil(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		g := ghost.New(t)

		var v interface{}

		result := ghost.BeNil(v)
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal("v is nil", result.Message))

		result = ghost.BeNil(nil)
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal("nil is nil", result.Message))
	})

	t.Run("non-nil", func(t *testing.T) {
		g := ghost.New(t)

		var v int

		result := ghost.BeNil(v)
		g.Should(ghost.BeFalse(result.Ok))
		g.Should(ghost.Equal("v is 0, not nil", result.Message))

		result = ghost.BeNil(-1 + 1)
		g.Should(ghost.BeFalse(result.Ok))
		g.Should(ghost.Equal("-1 + 1 is 0, not nil", result.Message))
	})
}

func TestBeTrue(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		g := ghost.New(t)

		v := true
		result := ghost.BeTrue(v)
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal("v is true", result.Message))

		result = ghost.BeTrue(true)
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal("true is true", result.Message))
	})

	t.Run("false", func(t *testing.T) {
		g := ghost.New(t)

		v := false
		result := ghost.BeTrue(v)
		g.Should(ghost.BeFalse(result.Ok))
		g.Should(ghost.Equal("v is false", result.Message))

		result = ghost.BeTrue(false)
		g.Should(ghost.BeFalse(result.Ok))
		g.Should(ghost.Equal("false is false", result.Message))
	})
}

func TestBeFalse(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		g := ghost.New(t)

		v := true
		result := ghost.BeFalse(v)
		g.Should(ghost.BeFalse(result.Ok))
		g.Should(ghost.Equal("v is true", result.Message))

		result = ghost.BeFalse(true)
		g.Should(ghost.BeFalse(result.Ok))
		g.Should(ghost.Equal("true is true", result.Message))
	})

	t.Run("false", func(t *testing.T) {
		g := ghost.New(t)

		v := false
		result := ghost.BeFalse(v)
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal("v is false", result.Message))

		result = ghost.BeFalse(false)
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal("false is false", result.Message))
	})
}

func TestBeZero(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		g := ghost.New(t)

		var v int
		result := ghost.BeZero(v)
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal("v is the zero value", result.Message))

		result = ghost.BeZero(0)
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal("0 is the zero value", result.Message))
	})

	t.Run("non-zero", func(t *testing.T) {
		g := ghost.New(t)

		v := 1
		result := ghost.BeZero(v)
		g.Should(ghost.BeFalse(result.Ok))
		g.Should(ghost.Equal("v is non-zero\nvalue: 1", result.Message))

		result = ghost.BeZero(1)
		g.Should(ghost.BeFalse(result.Ok))
		g.Should(ghost.Equal("1 is non-zero", result.Message))
	})
}

func TestContain(t *testing.T) {
	t.Run("contains <= 3", func(t *testing.T) {
		g := ghost.New(t)

		slice := []int{1, 2, 3}
		elem := 2

		result := ghost.Contain(slice, elem)
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal(`slice contains elem
slice:   [1 2 3]
element: 2
`, result.Message))

		result = ghost.Contain([]int{1, 2, 3}, 2)
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal(`[]int{1, 2, 3} contains 2
slice:   [1 2 3]
element: 2
`, result.Message))
	})

	t.Run("contains > 3", func(t *testing.T) {
		g := ghost.New(t)

		slice := []int{1, 2, 3, 4}
		elem := 2

		result := ghost.Contain(slice, elem)
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal(`slice contains elem
slice:   [
	1
>	2
	3
	4
]
element: 2
`, result.Message))

		result = ghost.Contain([]int{1, 2, 3, 4}, 2)
		g.Should(ghost.BeTrue(result.Ok))
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

		result := ghost.Contain(slice, elem)
		g.Should(ghost.BeFalse(result.Ok))
		g.Should(ghost.Equal(`slice does not contain elem
slice:   [1 2 3]
element: 5
`, result.Message))

		result = ghost.Contain([]int{1, 2, 3}, 5)
		g.Should(ghost.BeFalse(result.Ok))
		g.Should(ghost.Equal(`[]int{1, 2, 3} does not contain 5
slice:   [1 2 3]
element: 5
`, result.Message))
	})

	t.Run("does not contain > 3", func(t *testing.T) {
		g := ghost.New(t)

		slice := []int{1, 2, 3, 4}
		elem := 5

		result := ghost.Contain(slice, elem)
		g.Should(ghost.BeFalse(result.Ok))
		g.Should(ghost.Equal(`slice does not contain elem
slice:   [
	1
	2
	3
	4
]
element: 5
`, result.Message))

		result = ghost.Contain([]int{1, 2, 3, 4}, 5)
		g.Should(ghost.BeFalse(result.Ok))
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

		result := ghost.ContainString(outer, inner)
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal(`outer contains inner
str:    "foobar"
substr: "oob"
`, result.Message))

		result = ghost.ContainString("foobar", "oob")
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal(`"foobar" contains "oob"
str:    "foobar"
substr: "oob"
`, result.Message))
	})

	t.Run("does not contain", func(t *testing.T) {
		g := ghost.New(t)

		outer := "foobar"
		inner := "boo"

		result := ghost.ContainString(outer, inner)
		g.Should(ghost.BeFalse(result.Ok))
		g.Should(ghost.Equal(`outer does not contain inner
str:    "foobar"
substr: "boo"
`, result.Message))

		result = ghost.ContainString("foobar", "boo")
		g.Should(ghost.BeFalse(result.Ok))
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

		result := ghost.ContainString(outer, "two")
		g.Should(ghost.BeTrue(result.Ok))
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

		result := ghost.DeepEqual(want, got)
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal(`want == got
value: {foo [1 2]}
`, result.Message))

		result = ghost.DeepEqual(T{"foo", []int{1}}, T{"foo", []int{1}})
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal(`T{"foo", []int{1}} == T{"foo", []int{1}}
value: {foo [1]}
`, result.Message))
	})

	t.Run("unequal", func(t *testing.T) {
		g := ghost.New(t)

		type T struct {
			A string
			B int
		}

		want := T{"foo", 1}
		got := T{"bar", 0}

		result := ghost.DeepEqual(want, got)
		g.Should(ghost.BeFalse(result.Ok))

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

		result = ghost.DeepEqual(T{"foo", 1}, T{"bar", 0})
		g.Should(ghost.BeFalse(result.Ok))

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
	t.Run("equal", func(t *testing.T) {
		g := ghost.New(t)

		type T struct {
			A string
			B int
		}

		want := T{"foo", 1}
		got := T{"foo", 1}

		result := ghost.Equal(want, got)
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal(`want == got
value: {foo 1}
`, result.Message))

		result = ghost.Equal(T{"foo", 1}, T{"foo", 1})
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal(`T{"foo", 1} == T{"foo", 1}
value: {foo 1}
`, result.Message))
	})

	t.Run("unequal int", func(t *testing.T) {
		g := ghost.New(t)

		want := 1
		got := 0

		result := ghost.Equal(want, got)
		g.Should(ghost.BeFalse(result.Ok))

		wantText := `want != got
want: 1
got:  0
`
		g.Should(ghost.Equal(wantText, result.Message))

		result = ghost.Equal(1, 0)
		g.Should(ghost.BeFalse(result.Ok))

		wantText = `1 != 0
want: 1
got:  0
`
		g.Should(ghost.Equal(wantText, result.Message))
	})

	t.Run("unequal string short", func(t *testing.T) {
		g := ghost.New(t)

		want := "foo"
		got := "bar"

		result := ghost.Equal(want, got)
		g.Should(ghost.BeFalse(result.Ok))

		wantText := `want != got
want: "foo"
got:  "bar"
`
		g.Should(ghost.Equal(wantText, result.Message))

		result = ghost.Equal("foo", "bar")
		g.Should(ghost.BeFalse(result.Ok))

		wantText = `"foo" != "bar"
want: "foo"
got:  "bar"
`
		g.Should(ghost.Equal(wantText, result.Message))
	})

	t.Run("unequal string long", func(t *testing.T) {
		g := ghost.New(t)

		want := "foo\nbar\nbaz"
		got := "bar"

		result := ghost.Equal(want, got)
		g.Should(ghost.BeFalse(result.Ok))

		wantText := `want != got
want: ` + `
"""
foo
bar
baz
"""

got:  "bar"
`
		g.Should(ghost.Equal(wantText, result.Message))

		result = ghost.Equal("foo\nbar\nbaz", "bar")
		g.Should(ghost.BeFalse(result.Ok))

		wantText = `"foo\nbar\nbaz" != "bar"
want: ` + `
"""
foo
bar
baz
"""

got:  "bar"
`
		g.Should(ghost.Equal(wantText, result.Message))
	})

	t.Run("unequal struct", func(t *testing.T) {
		g := ghost.New(t)

		type T struct {
			A string
			B int
		}

		want := T{"foo", 1}
		got := T{"bar", 0}

		result := ghost.Equal(want, got)
		g.Should(ghost.BeFalse(result.Ok))

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

		result = ghost.Equal(T{"foo", 1}, T{"bar", 0})
		g.Should(ghost.BeFalse(result.Ok))

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

func TestError(t *testing.T) {
	t.Run("non-nil", func(t *testing.T) {
		g := ghost.New(t)

		err := errors.New("oopsie")

		result := ghost.Error(err)
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal(`err has error value: oopsie`, result.Message))

		result = ghost.Error(errors.New("oopsie"))
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal(`errors.New("oopsie") has error value: oopsie`, result.Message))
	})

	t.Run("nil", func(t *testing.T) {
		g := ghost.New(t)

		var err error

		result := ghost.Error(err)
		g.Should(ghost.BeFalse(result.Ok))
		g.Should(ghost.Equal(`err is nil`, result.Message))

		result = ghost.Error(nil)
		g.Should(ghost.BeFalse(result.Ok))
		g.Should(ghost.Equal(`nil is nil`, result.Message))
	})
}

func TestErrorContaining(t *testing.T) {
	t.Run("contains", func(t *testing.T) {
		g := ghost.New(t)

		err := errors.New("foobar")
		msg := "oob"

		result := ghost.ErrorContaining(msg, err)
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal(
			`err contains error message "oob": foobar`,
			result.Message,
		))

		result = ghost.ErrorContaining("oob", errors.New("foobar"))
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal(
			`errors.New("foobar") contains error message "oob": foobar`,
			result.Message,
		))
	})

	t.Run("does not contain", func(t *testing.T) {
		g := ghost.New(t)

		err := errors.New("foobar")
		msg := "boo"

		result := ghost.ErrorContaining(msg, err)
		g.Should(ghost.BeFalse(result.Ok))
		g.Should(ghost.Equal(
			`err does not contain error message "boo": foobar`,
			result.Message,
		))

		result = ghost.ErrorContaining("boo", errors.New("foobar"))
		g.Should(ghost.BeFalse(result.Ok))
		g.Should(ghost.Equal(
			`errors.New("foobar") does not contain error message "boo": foobar`,
			result.Message,
		))
	})

	t.Run("nil", func(t *testing.T) {
		g := ghost.New(t)

		var err error
		msg := "boo"

		result := ghost.ErrorContaining(msg, err)
		g.Should(ghost.BeFalse(result.Ok))
		g.Should(ghost.Equal(`err is nil; missing error message msg: boo`, result.Message))

		result = ghost.ErrorContaining("boo", nil)
		g.Should(ghost.BeFalse(result.Ok))
		g.Should(ghost.Equal(`nil is nil; missing error message: boo`, result.Message))
	})
}

func TestErrorEqual(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		g := ghost.New(t)

		err := errors.New("foobar")
		msg := "foobar"

		result := ghost.ErrorEqual(msg, err)
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal(
			`err equals error message "foobar": foobar`,
			result.Message,
		))

		result = ghost.ErrorEqual("foobar", errors.New("foobar"))
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal(
			`errors.New("foobar") equals error message "foobar": foobar`,
			result.Message,
		))
	})

	t.Run("not equal", func(t *testing.T) {
		g := ghost.New(t)

		err := errors.New("foobar")
		msg := "boo"

		result := ghost.ErrorEqual(msg, err)
		g.Should(ghost.BeFalse(result.Ok))
		g.Should(ghost.Equal(
			`err does not equal error message "boo": foobar`,
			result.Message,
		))

		result = ghost.ErrorEqual("boo", errors.New("foobar"))
		g.Should(ghost.BeFalse(result.Ok))
		g.Should(ghost.Equal(
			`errors.New("foobar") does not equal error message "boo": foobar`,
			result.Message,
		))
	})

	t.Run("nil", func(t *testing.T) {
		g := ghost.New(t)

		var err error
		msg := "boo"

		result := ghost.ErrorEqual(msg, err)
		g.Should(ghost.BeFalse(result.Ok))
		g.Should(ghost.Equal(`err is nil; want message: boo`, result.Message))

		result = ghost.ErrorEqual("boo", nil)
		g.Should(ghost.BeFalse(result.Ok))
		g.Should(ghost.Equal(`nil is nil; want message: boo`, result.Message))
	})
}

func TestJSONEqual(t *testing.T) {
	// TODO: Write me
	_ = t
}

func TestLen(t *testing.T) {
	t.Run("equal <= 3", func(t *testing.T) {
		g := ghost.New(t)

		wantLen := 3
		slice := []string{"a", "b", "c"}

		result := ghost.Len(wantLen, slice)
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal(`want slice length 3, got 3
slice: [a b c]
`, result.Message))

		result = ghost.Len(3, []string{"a", "b", "c"})
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal(`want []string{"a", "b", "c"} length 3, got 3
slice: [a b c]
`, result.Message))
	})

	t.Run("equal > 3", func(t *testing.T) {
		g := ghost.New(t)

		wantLen := 4
		slice := []string{"a", "b", "c", "d"}

		result := ghost.Len(wantLen, slice)
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal(`want slice length 4, got 4
slice: [
	a
	b
	c
	d
]
`, result.Message))

		result = ghost.Len(4, []string{"a", "b", "c", "d"})
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal(`want []string{"a", "b", "c", "d"} length 4, got 4
slice: [
	a
	b
	c
	d
]
`, result.Message))
	})
}

func TestPanic(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		g := ghost.New(t)

		f := func() { panic(errors.New("oh no")) }

		result := ghost.Panic(f)
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal("function f panicked with value: oh no", result.Message))

		result = ghost.Panic(func() { panic(errors.New("oh no")) })
		g.Should(ghost.BeTrue(result.Ok))
		g.Should(ghost.Equal(`function panicked with value: oh no
func() {
	panic(errors.New("oh no"))
}
`, result.Message))
	})

	t.Run("no panic", func(t *testing.T) {
		g := ghost.New(t)

		f := func() {}

		result := ghost.Panic(f)
		g.Should(ghost.BeFalse(result.Ok))
		g.Should(ghost.Equal("function f did not panic", result.Message))

		result = ghost.Panic(func() {})
		g.Should(ghost.BeFalse(result.Ok))
		g.Should(ghost.Equal(`function did not panic
func() {
}
`, result.Message))
	})
}
