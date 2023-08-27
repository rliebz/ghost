package be_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/rliebz/ghost"
	"github.com/rliebz/ghost/be"
)

func TestClose(t *testing.T) {
	t.Run("in delta", func(t *testing.T) {
		g := ghost.New(t)

		want := 32.5
		got := 32.0

		result := be.Close(want, got, 1)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(
			"delta 0.5 between want (32.5) and got (32) is within 1",
			result.Message,
		))

		result = be.Close(32.5, 32.0, 1.0)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(
			"delta 0.5 between 32.5 and 32.0 is within 1",
			result.Message,
		))

		result = be.Close(32.0, 32.5, 1.0)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(
			"delta 0.5 between 32.0 and 32.5 is within 1",
			result.Message,
		))
	})

	t.Run("not in delta", func(t *testing.T) {
		g := ghost.New(t)

		want := 32.5
		got := 32.0

		result := be.Close(want, got, 0.3)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(
			"delta 0.5 between want (32.5) and got (32) is not within 0.3",
			result.Message,
		))

		result = be.Close(32.5, 32.0, 0.3)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(
			"delta 0.5 between 32.5 and 32.0 is not within 0.3",
			result.Message,
		))
	})
}

func TestDeepEqual(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		g := ghost.New(t)

		type T struct {
			A string
			b []int
		}

		want := T{"foo", []int{1, 2}}
		got := T{"foo", []int{1, 2}}

		result := be.DeepEqual(want, got)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(`want == got
value: {foo [1 2]}
`, result.Message))

		result = be.DeepEqual(T{"foo", []int{1}}, T{"foo", []int{1}})
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(`T{"foo", []int{1}} == T{"foo", []int{1}}
value: {foo [1]}
`, result.Message))
	})

	t.Run("unequal", func(t *testing.T) {
		g := ghost.New(t)

		type T struct {
			A string
			b int
		}

		want := T{"foo", 1}
		got := T{"bar", 0}

		result := be.DeepEqual(want, got)
		g.Should(be.False(result.Ok))

		// Keep the diff small, because we don't want to test cmp.Diff
		wantText := `want != got
diff (-want +got):
  be_test.T{
- 	A: "foo",
+ 	A: "bar",
- 	b: 1,
+ 	b: 0,
  }

`
		result.Message = strings.ReplaceAll(result.Message, "\u00a0", " ")
		g.Should(be.Equal(wantText, result.Message))

		result = be.DeepEqual(T{"foo", 1}, T{"bar", 0})
		g.Should(be.False(result.Ok))

		wantText = `T{"foo", 1} != T{"bar", 0}
diff (-want +got):
  be_test.T{
- 	A: "foo",
+ 	A: "bar",
- 	b: 1,
+ 	b: 0,
  }

`
		result.Message = strings.ReplaceAll(result.Message, "\u00a0", " ")
		g.Should(be.Equal(wantText, result.Message))
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

		result := be.Equal(want, got)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(`want == got
value: {foo 1}
`, result.Message))

		result = be.Equal(T{"foo", 1}, T{"foo", 1})
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(`T{"foo", 1} == T{"foo", 1}
value: {foo 1}
`, result.Message))
	})

	t.Run("equal simple", func(t *testing.T) {
		g := ghost.New(t)

		got := 3

		result := be.Equal(3, got)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(`3 == got`, result.Message))
	})

	t.Run("unequal int", func(t *testing.T) {
		g := ghost.New(t)

		want := 1
		got := 0

		result := be.Equal(want, got)
		g.Should(be.False(result.Ok))

		wantText := `want != got
want: 1
got:  0
`
		g.Should(be.Equal(wantText, result.Message))

		result = be.Equal(1, 0)
		g.Should(be.False(result.Ok))

		wantText = `1 != 0
want: 1
got:  0
`
		g.Should(be.Equal(wantText, result.Message))
	})

	t.Run("unequal string short", func(t *testing.T) {
		g := ghost.New(t)

		want := "foo"
		got := "bar"

		result := be.Equal(want, got)
		g.Should(be.False(result.Ok))

		wantText := `want != got
want: "foo"
got:  "bar"
`
		g.Should(be.Equal(wantText, result.Message))

		result = be.Equal("foo", "bar")
		g.Should(be.False(result.Ok))

		wantText = `"foo" != "bar"
want: "foo"
got:  "bar"
`
		g.Should(be.Equal(wantText, result.Message))
	})

	t.Run("unequal string long", func(t *testing.T) {
		g := ghost.New(t)

		want := "foo\nbar\nbaz"
		got := "bar"

		result := be.Equal(want, got)
		g.Should(be.False(result.Ok))

		wantText := `want != got
want: ` + `
"""
foo
bar
baz
"""

got:  "bar"
`
		g.Should(be.Equal(wantText, result.Message))

		result = be.Equal("foo\nbar\nbaz", "bar")
		g.Should(be.False(result.Ok))

		wantText = `"foo\nbar\nbaz" != "bar"
want: ` + `
"""
foo
bar
baz
"""

got:  "bar"
`
		g.Should(be.Equal(wantText, result.Message))
	})

	t.Run("unequal struct", func(t *testing.T) {
		g := ghost.New(t)

		type T struct {
			A string
			B int
		}

		want := T{"foo", 1}
		got := T{"bar", 0}

		result := be.Equal(want, got)
		g.Should(be.False(result.Ok))

		// Keep the diff small, because we don't want to test cmp.Diff
		wantText := `want != got
diff (-want +got):
  be_test.T{
- 	A: "foo",
+ 	A: "bar",
- 	B: 1,
+ 	B: 0,
  }

`
		result.Message = strings.ReplaceAll(result.Message, "\u00a0", " ")
		g.Should(be.Equal(wantText, result.Message))

		result = be.Equal(T{"foo", 1}, T{"bar", 0})
		g.Should(be.False(result.Ok))

		wantText = `T{"foo", 1} != T{"bar", 0}
diff (-want +got):
  be_test.T{
- 	A: "foo",
+ 	A: "bar",
- 	B: 1,
+ 	B: 0,
  }

`
		result.Message = strings.ReplaceAll(result.Message, "\u00a0", " ")
		g.Should(be.Equal(wantText, result.Message))
	})
}

func TestError(t *testing.T) {
	t.Run("non-nil", func(t *testing.T) {
		g := ghost.New(t)

		err := errors.New("oopsie")

		result := be.Error(err)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(`err has error value: oopsie`, result.Message))

		result = be.Error(errors.New("oopsie"))
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(`errors.New("oopsie") has error value: oopsie`, result.Message))
	})

	t.Run("nil", func(t *testing.T) {
		g := ghost.New(t)

		var err error

		result := be.Error(err)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(`err is nil`, result.Message))

		result = be.Error(nil)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(`nil is nil`, result.Message))
	})
}

func TestErrorContaining(t *testing.T) {
	t.Run("contains", func(t *testing.T) {
		g := ghost.New(t)

		err := errors.New("foobar")
		msg := "oob"

		result := be.ErrorContaining(msg, err)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(
			`err contains error message "oob": foobar`,
			result.Message,
		))

		result = be.ErrorContaining("oob", errors.New("foobar"))
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(
			`errors.New("foobar") contains error message "oob": foobar`,
			result.Message,
		))
	})

	t.Run("does not contain", func(t *testing.T) {
		g := ghost.New(t)

		err := errors.New("foobar")
		msg := "boo"

		result := be.ErrorContaining(msg, err)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(
			`err does not contain error message "boo": foobar`,
			result.Message,
		))

		result = be.ErrorContaining("boo", errors.New("foobar"))
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(
			`errors.New("foobar") does not contain error message "boo": foobar`,
			result.Message,
		))
	})

	t.Run("nil", func(t *testing.T) {
		g := ghost.New(t)

		var err error
		msg := "boo"

		result := be.ErrorContaining(msg, err)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(`err is nil; missing error message msg: boo`, result.Message))

		result = be.ErrorContaining("boo", nil)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(`nil is nil; missing error message: boo`, result.Message))
	})
}

func TestErrorEqual(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		g := ghost.New(t)

		err := errors.New("foobar")
		msg := "foobar"

		result := be.ErrorEqual(msg, err)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(
			`err equals error message "foobar": foobar`,
			result.Message,
		))

		result = be.ErrorEqual("foobar", errors.New("foobar"))
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(
			`errors.New("foobar") equals error message "foobar": foobar`,
			result.Message,
		))
	})

	t.Run("not equal", func(t *testing.T) {
		g := ghost.New(t)

		err := errors.New("foobar")
		msg := "boo"

		result := be.ErrorEqual(msg, err)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(
			`err does not equal error message "boo": foobar`,
			result.Message,
		))

		result = be.ErrorEqual("boo", errors.New("foobar"))
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(
			`errors.New("foobar") does not equal error message "boo": foobar`,
			result.Message,
		))
	})

	t.Run("nil", func(t *testing.T) {
		g := ghost.New(t)

		var err error
		msg := "boo"

		result := be.ErrorEqual(msg, err)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(`err is nil; want message: boo`, result.Message))

		result = be.ErrorEqual("boo", nil)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(`nil is nil; want message: boo`, result.Message))
	})
}

func TestFalse(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		g := ghost.New(t)

		v := true
		result := be.False(v)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal("v is true", result.Message))

		result = be.False(true)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal("true is true", result.Message))
	})

	t.Run("false", func(t *testing.T) {
		g := ghost.New(t)

		v := false
		result := be.False(v)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal("v is false", result.Message))

		result = be.False(false)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal("false is false", result.Message))
	})
}

func TestJSONEqual(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		g := ghost.New(t)

		want := `{"foo": "value", "bar": [1, 2]}`
		got := `{"bar": [1, 2], "foo": "value"}`

		result := be.JSONEqual(want, got)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal("want and got are JSON equal", result.Message))

		result = be.JSONEqual(`{"foo": "value", "bar": [1, 2]}`, `{"bar": [1, 2], "foo": "value"}`)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(
			"`{\"foo\": \"value\", \"bar\": [1, 2]}` and "+
				"`{\"bar\": [1, 2], \"foo\": \"value\"}` are JSON equal",
			result.Message,
		))
	})

	t.Run("not equal", func(t *testing.T) {
		g := ghost.New(t)

		want := `{"foo": "value", "bar": [1, 2]}`
		got := `{"bar": [2, 1], "foo": "other"}`

		result := be.JSONEqual(want, got)
		g.Should(be.False(result.Ok))
		g.Should(be.StringContaining("want and got are not JSON equal", result.Message))

		result = be.JSONEqual(`{"foo": "value", "bar": [1, 2]}`, `{"bar": [2, 1], "foo": "other"}`)
		g.Should(be.False(result.Ok))
		g.Should(be.StringContaining(
			"`{\"foo\": \"value\", \"bar\": [1, 2]}` and "+
				"`{\"bar\": [2, 1], \"foo\": \"other\"}` are not JSON equal",
			result.Message,
		))
	})

	t.Run("invalid json", func(t *testing.T) {
		g := ghost.New(t)

		valid := `{"foo": "value", "bar": [1, 2]}`
		invalid := `{{`
		invalid2 := `{{`

		result := be.JSONEqual(valid, invalid)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(`invalid is not valid JSON
value: {{`, result.Message))

		result = be.JSONEqual(invalid, valid)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(`invalid is not valid JSON
value: {{`, result.Message))

		result = be.JSONEqual(invalid, invalid2)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(`invalid and invalid2 are not valid JSON
want:
{{

got:
{{`, result.Message))

		result = be.JSONEqual(`{"foo": "value", "bar": [1, 2]}`, `{"bar": [1, 2], "foo": "value"}`)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(
			"`{\"foo\": \"value\", \"bar\": [1, 2]}` and "+
				"`{\"bar\": [1, 2], \"foo\": \"value\"}` are JSON equal",
			result.Message,
		))
	})
}

func TestMapLen(t *testing.T) {
	t.Run("equal <= 3", func(t *testing.T) {
		g := ghost.New(t)

		wantLen := 3
		m := map[string]int{"a": 1, "b": 2, "c": 3}

		result := be.MapLen(wantLen, m)
		g.Should(be.True(result.Ok))
		g.Should(be.StringContaining(`want m length 3, got 3`, result.Message))

		result = be.MapLen(3, map[string]int{"a": 1, "b": 2, "c": 3})
		g.Should(be.True(result.Ok))
		g.Should(be.StringContaining(
			`want map[string]int{"a": 1, "b": 2, "c": 3} length 3, got 3`,
			result.Message,
		))
	})

	t.Run("equal > 3", func(t *testing.T) {
		g := ghost.New(t)

		wantLen := 4
		m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}

		result := be.MapLen(wantLen, m)
		g.Should(be.True(result.Ok))
		g.Should(be.StringContaining(`want m length 4, got 4`, result.Message))

		result = be.MapLen(4, map[string]int{"a": 1, "b": 2, "c": 3, "d": 4})
		g.Should(be.True(result.Ok))
		g.Should(be.StringContaining(
			`want map[string]int{"a": 1, "b": 2, "c": 3, "d": 4} length 4, got 4`,
			result.Message,
		))
	})

	t.Run("not equal <= 3", func(t *testing.T) {
		g := ghost.New(t)

		wantLen := 2
		m := map[string]int{"a": 1, "b": 2, "c": 3}

		result := be.MapLen(wantLen, m)
		g.Should(be.False(result.Ok))
		g.Should(be.StringContaining(`want m length 2, got 3`, result.Message))

		result = be.MapLen(2, map[string]int{"a": 1, "b": 2, "c": 3})
		g.Should(be.False(result.Ok))
		g.Should(be.StringContaining(
			`want map[string]int{"a": 1, "b": 2, "c": 3} length 2, got 3`,
			result.Message,
		))
	})

	t.Run("not equal > 3", func(t *testing.T) {
		g := ghost.New(t)

		wantLen := 3
		m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}

		result := be.MapLen(wantLen, m)
		g.Should(be.False(result.Ok))
		g.Should(be.StringContaining(`want m length 3, got 4`, result.Message))

		result = be.MapLen(3, map[string]int{"a": 1, "b": 2, "c": 3, "d": 4})
		g.Should(be.False(result.Ok))
		g.Should(be.StringContaining(
			`want map[string]int{"a": 1, "b": 2, "c": 3, "d": 4} length 3, got 4`,
			result.Message,
		))
	})
}

func TestNil(t *testing.T) {
	t.Run("untyped nil", func(t *testing.T) {
		g := ghost.New(t)

		var v interface{}

		result := be.Nil(v)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal("v is nil", result.Message))

		result = be.Nil(nil)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal("nil is nil", result.Message))
	})

	t.Run("typed nil", func(t *testing.T) {
		g := ghost.New(t)

		var v *int
		var i interface{} = v

		result := be.Nil(i)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal("i is nil", result.Message))

		result = be.Nil((*int)(nil))
		g.Should(be.True(result.Ok))
		g.Should(be.Equal("(*int)(nil) is nil", result.Message))
	})

	t.Run("non-nil", func(t *testing.T) {
		g := ghost.New(t)

		var v int

		result := be.Nil(v)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal("v is 0, not nil", result.Message))

		result = be.Nil(-1 + 1)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal("-1 + 1 is 0, not nil", result.Message))
	})
}

func TestPanic(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		g := ghost.New(t)

		f := func() { panic(errors.New("oh no")) }

		result := be.Panic(f)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal("function f panicked with value: oh no", result.Message))

		result = be.Panic(func() { panic(errors.New("oh no")) })
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(`function panicked with value: oh no
func() {
	panic(errors.New("oh no"))
}
`, result.Message))
	})

	t.Run("no panic", func(t *testing.T) {
		g := ghost.New(t)

		f := func() {}

		result := be.Panic(f)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal("function f did not panic", result.Message))

		result = be.Panic(func() {})
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(`function did not panic
func() {
}
`, result.Message))
	})
}

func TestSliceContaining(t *testing.T) {
	t.Run("contains <= 3", func(t *testing.T) {
		g := ghost.New(t)

		slice := []int{1, 2, 3}
		elem := 2

		result := be.SliceContaining(elem, slice)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(`slice contains elem
element: 2
slice:   [1 2 3]
`, result.Message))

		result = be.SliceContaining(2, []int{1, 2, 3})
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(`[]int{1, 2, 3} contains 2
element: 2
slice:   [1 2 3]
`, result.Message))
	})

	t.Run("contains > 3", func(t *testing.T) {
		g := ghost.New(t)

		slice := []int{1, 2, 3, 4}
		elem := 2

		result := be.SliceContaining(elem, slice)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(`slice contains elem
element: 2
slice:   [
	1
>	2
	3
	4
]
`, result.Message))

		result = be.SliceContaining(2, []int{1, 2, 3, 4})
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(`[]int{1, 2, 3, 4} contains 2
element: 2
slice:   [
	1
>	2
	3
	4
]
`, result.Message))
	})

	t.Run("does not contain <= 3", func(t *testing.T) {
		g := ghost.New(t)

		slice := []int{1, 2, 3}
		elem := 5

		result := be.SliceContaining(elem, slice)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(`slice does not contain elem
element: 5
slice:   [1 2 3]
`, result.Message))

		result = be.SliceContaining(5, []int{1, 2, 3})
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(`[]int{1, 2, 3} does not contain 5
element: 5
slice:   [1 2 3]
`, result.Message))
	})

	t.Run("does not contain > 3", func(t *testing.T) {
		g := ghost.New(t)

		slice := []int{1, 2, 3, 4}
		elem := 5

		result := be.SliceContaining(elem, slice)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(`slice does not contain elem
element: 5
slice:   [
	1
	2
	3
	4
]
`, result.Message))

		result = be.SliceContaining(5, []int{1, 2, 3, 4})
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(`[]int{1, 2, 3, 4} does not contain 5
element: 5
slice:   [
	1
	2
	3
	4
]
`, result.Message))
	})
}

func TestSliceLen(t *testing.T) {
	t.Run("equal <= 3", func(t *testing.T) {
		g := ghost.New(t)

		wantLen := 3
		slice := []string{"a", "b", "c"}

		result := be.SliceLen(wantLen, slice)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(`want slice length 3, got 3
slice: [a b c]
`, result.Message))

		result = be.SliceLen(3, []string{"a", "b", "c"})
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(`want []string{"a", "b", "c"} length 3, got 3
slice: [a b c]
`, result.Message))
	})

	t.Run("equal > 3", func(t *testing.T) {
		g := ghost.New(t)

		wantLen := 4
		slice := []string{"a", "b", "c", "d"}

		result := be.SliceLen(wantLen, slice)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(`want slice length 4, got 4
slice: [
	a
	b
	c
	d
]
`, result.Message))

		result = be.SliceLen(4, []string{"a", "b", "c", "d"})
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(`want []string{"a", "b", "c", "d"} length 4, got 4
slice: [
	a
	b
	c
	d
]
`, result.Message))
	})

	t.Run("not equal <= 3", func(t *testing.T) {
		g := ghost.New(t)

		wantLen := 2
		slice := []string{"a", "b", "c"}

		result := be.SliceLen(wantLen, slice)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(`want slice length 2, got 3
slice: [a b c]
`, result.Message))

		result = be.SliceLen(2, []string{"a", "b", "c"})
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(`want []string{"a", "b", "c"} length 2, got 3
slice: [a b c]
`, result.Message))
	})

	t.Run("not equal > 3", func(t *testing.T) {
		g := ghost.New(t)

		wantLen := 3
		slice := []string{"a", "b", "c", "d"}

		result := be.SliceLen(wantLen, slice)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(`want slice length 3, got 4
slice: [
	a
	b
	c
	d
]
`, result.Message))

		result = be.SliceLen(3, []string{"a", "b", "c", "d"})
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(`want []string{"a", "b", "c", "d"} length 3, got 4
slice: [
	a
	b
	c
	d
]
`, result.Message))
	})
}

func TestStringContaining(t *testing.T) {
	t.Run("contains", func(t *testing.T) {
		g := ghost.New(t)

		outer := "foobar"
		inner := "oob"

		result := be.StringContaining(inner, outer)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(`outer contains inner
substr: "oob"
str:    "foobar"
`, result.Message))

		result = be.StringContaining("oob", "foobar")
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(`"foobar" contains "oob"
substr: "oob"
str:    "foobar"
`, result.Message))
	})

	t.Run("does not contain", func(t *testing.T) {
		g := ghost.New(t)

		outer := "foobar"
		inner := "boo"

		result := be.StringContaining(inner, outer)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(`outer does not contain inner
substr: "boo"
str:    "foobar"
`, result.Message))

		result = be.StringContaining("boo", "foobar")
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(`"foobar" does not contain "boo"
substr: "boo"
str:    "foobar"
`, result.Message))
	})

	t.Run("multiline", func(t *testing.T) {
		g := ghost.New(t)

		outer := `one
two
three
`

		result := be.StringContaining("two", outer)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(`outer contains "two"
substr: "two"
str:    `+`
"""
one
two
three

"""

`, result.Message))
	})
}

func TestTrue(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		g := ghost.New(t)

		v := true
		result := be.True(v)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal("v is true", result.Message))

		result = be.True(true)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal("true is true", result.Message))
	})

	t.Run("false", func(t *testing.T) {
		g := ghost.New(t)

		v := false
		result := be.True(v)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal("v is false", result.Message))

		result = be.True(false)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal("false is false", result.Message))
	})
}

func TestZero(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		g := ghost.New(t)

		var v int
		result := be.Zero(v)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal("v is the zero value", result.Message))

		result = be.Zero(0)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal("0 is the zero value", result.Message))
	})

	t.Run("non-zero", func(t *testing.T) {
		g := ghost.New(t)

		v := 1
		result := be.Zero(v)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal("v is non-zero\nvalue: 1", result.Message))

		result = be.Zero(1)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal("1 is non-zero", result.Message))
	})
}
