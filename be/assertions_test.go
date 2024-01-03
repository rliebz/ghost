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

		got := 32.0
		want := 32.5

		result := be.Close(got, want, 1)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(
			result.Message,
			"delta 0.5 between got (32) and want (32.5) is within 1",
		))

		result = be.Close(32.0, 32.5, 1.0)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(
			result.Message,
			"delta 0.5 between 32.0 and 32.5 is within 1",
		))

		result = be.Close(32.5, 32.0, 1.0)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(
			result.Message,
			"delta 0.5 between 32.5 and 32.0 is within 1",
		))
	})

	t.Run("not in delta", func(t *testing.T) {
		g := ghost.New(t)

		got := 32.0
		want := 32.5

		result := be.Close(got, want, 0.3)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(
			result.Message,
			"delta 0.5 between got (32) and want (32.5) is not within 0.3",
		))

		result = be.Close(32.0, 32.5, 0.3)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(
			result.Message,
			"delta 0.5 between 32.0 and 32.5 is not within 0.3",
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

		got := T{"foo", []int{1, 2}}
		want := T{"foo", []int{1, 2}}

		result := be.DeepEqual(got, want)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `got == want
value: {foo [1 2]}
`))

		result = be.DeepEqual(T{"foo", []int{1}}, T{"foo", []int{1}})
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `T{"foo", []int{1}} == T{"foo", []int{1}}
value: {foo [1]}
`))
	})

	t.Run("unequal", func(t *testing.T) {
		g := ghost.New(t)

		type T struct {
			A string
			b int
		}

		got := T{"bar", 0}
		want := T{"foo", 1}

		result := be.DeepEqual(got, want)
		g.Should(be.False(result.Ok))

		// Keep the diff small, because we don't want to test cmp.Diff
		wantText := `got != want
diff (-want +got):
  be_test.T{
- 	A: "foo",
+ 	A: "bar",
- 	b: 1,
+ 	b: 0,
  }

`
		result.Message = strings.ReplaceAll(result.Message, "\u00a0", " ")
		g.Should(be.Equal(result.Message, wantText))

		result = be.DeepEqual(T{"bar", 0}, T{"foo", 1})
		g.Should(be.False(result.Ok))

		wantText = `T{"bar", 0} != T{"foo", 1}
diff (-want +got):
  be_test.T{
- 	A: "foo",
+ 	A: "bar",
- 	b: 1,
+ 	b: 0,
  }

`
		result.Message = strings.ReplaceAll(result.Message, "\u00a0", " ")
		g.Should(be.Equal(result.Message, wantText))
	})
}

func TestEqual(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		g := ghost.New(t)

		type T struct {
			A string
			B int
		}

		got := T{"foo", 1}
		want := T{"foo", 1}

		result := be.Equal(got, want)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `got == want
value: {foo 1}
`))

		result = be.Equal(T{"foo", 1}, T{"foo", 1})
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `T{"foo", 1} == T{"foo", 1}
value: {foo 1}
`))
	})

	t.Run("equal simple", func(t *testing.T) {
		g := ghost.New(t)

		got := 3

		result := be.Equal(got, 3)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `got == 3`))
	})

	t.Run("unequal int", func(t *testing.T) {
		g := ghost.New(t)

		got := 0
		want := 1

		result := be.Equal(got, want)
		g.Should(be.False(result.Ok))

		wantText := `got != want
got:  0
want: 1
`
		g.Should(be.Equal(result.Message, wantText))

		result = be.Equal(0, 1)
		g.Should(be.False(result.Ok))

		wantText = `0 != 1
got:  0
want: 1
`
		g.Should(be.Equal(result.Message, wantText))
	})

	t.Run("unequal string short", func(t *testing.T) {
		g := ghost.New(t)

		got := "bar"
		want := "foo"

		result := be.Equal(got, want)
		g.Should(be.False(result.Ok))

		wantText := `got != want
got:  "bar"
want: "foo"
`
		g.Should(be.Equal(result.Message, wantText))

		result = be.Equal("bar", "foo")
		g.Should(be.False(result.Ok))

		wantText = `"bar" != "foo"
got:  "bar"
want: "foo"
`
		g.Should(be.Equal(result.Message, wantText))
	})

	t.Run("unequal string long", func(t *testing.T) {
		g := ghost.New(t)

		got := "bar"
		want := "foo\nbar\nbaz"

		result := be.Equal(got, want)
		g.Should(be.False(result.Ok))

		wantText := `got != want
got:  "bar"
want: ` + `
"""
foo
bar
baz
"""
`
		g.Should(be.Equal(result.Message, wantText))

		result = be.Equal("bar", "foo\nbar\nbaz")
		g.Should(be.False(result.Ok))

		wantText = `"bar" != "foo\nbar\nbaz"
got:  "bar"
want: ` + `
"""
foo
bar
baz
"""
`
		g.Should(be.Equal(result.Message, wantText))
	})

	t.Run("unequal struct", func(t *testing.T) {
		g := ghost.New(t)

		type T struct {
			A string
			B int
		}

		got := T{"bar", 0}
		want := T{"foo", 1}

		result := be.Equal(got, want)
		g.Should(be.False(result.Ok))

		// Keep the diff small, because we don't want to test cmp.Diff
		wantText := `got != want
diff (-want +got):
  be_test.T{
- 	A: "foo",
+ 	A: "bar",
- 	B: 1,
+ 	B: 0,
  }

`
		result.Message = strings.ReplaceAll(result.Message, "\u00a0", " ")
		g.Should(be.Equal(result.Message, wantText))

		result = be.Equal(T{"bar", 0}, T{"foo", 1})
		g.Should(be.False(result.Ok))

		wantText = `T{"bar", 0} != T{"foo", 1}
diff (-want +got):
  be_test.T{
- 	A: "foo",
+ 	A: "bar",
- 	B: 1,
+ 	B: 0,
  }

`
		result.Message = strings.ReplaceAll(result.Message, "\u00a0", " ")
		g.Should(be.Equal(result.Message, wantText))
	})
}

func TestError(t *testing.T) {
	t.Run("non-nil", func(t *testing.T) {
		g := ghost.New(t)

		err := errors.New("oopsie")

		result := be.Error(err)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `err has error value: oopsie`))

		result = be.Error(errors.New("oopsie"))
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `errors.New("oopsie") has error value: oopsie`))
	})

	t.Run("nil", func(t *testing.T) {
		g := ghost.New(t)

		var err error

		result := be.Error(err)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `err is nil`))

		result = be.Error(nil)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `nil is nil`))
	})
}

func TestErrorContaining(t *testing.T) {
	t.Run("contains", func(t *testing.T) {
		g := ghost.New(t)

		err := errors.New("foobar")
		msg := "oob"

		result := be.ErrorContaining(err, msg)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(
			result.Message,
			`err contains error message "oob": foobar`,
		))

		result = be.ErrorContaining(errors.New("foobar"), "oob")
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(
			result.Message,
			`errors.New("foobar") contains error message "oob": foobar`,
		))
	})

	t.Run("does not contain", func(t *testing.T) {
		g := ghost.New(t)

		err := errors.New("foobar")
		msg := "boo"

		result := be.ErrorContaining(err, msg)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(
			result.Message,
			`err does not contain error message "boo": foobar`,
		))

		result = be.ErrorContaining(errors.New("foobar"), "boo")
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(
			result.Message,
			`errors.New("foobar") does not contain error message "boo": foobar`,
		))
	})

	t.Run("nil", func(t *testing.T) {
		g := ghost.New(t)

		var err error
		msg := "boo"

		result := be.ErrorContaining(err, msg)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `err is nil; missing error message msg: boo`))

		result = be.ErrorContaining(nil, "boo")
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `nil is nil; missing error message: boo`))
	})
}

func TestErrorEqual(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		g := ghost.New(t)

		err := errors.New("foobar")
		msg := "foobar"

		result := be.ErrorEqual(err, msg)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(
			result.Message,
			`err equals error message "foobar": foobar`,
		))

		result = be.ErrorEqual(errors.New("foobar"), "foobar")
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(
			result.Message,
			`errors.New("foobar") equals error message "foobar": foobar`,
		))
	})

	t.Run("not equal", func(t *testing.T) {
		g := ghost.New(t)

		err := errors.New("foobar")
		msg := "boo"

		result := be.ErrorEqual(err, msg)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(
			result.Message,
			`err does not equal error message "boo": foobar`,
		))

		result = be.ErrorEqual(errors.New("foobar"), "boo")
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(
			result.Message,
			`errors.New("foobar") does not equal error message "boo": foobar`,
		))
	})

	t.Run("nil", func(t *testing.T) {
		g := ghost.New(t)

		var err error
		msg := "boo"

		result := be.ErrorEqual(err, msg)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `err is nil; want message: boo`))

		result = be.ErrorEqual(nil, "boo")
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `nil is nil; want message: boo`))
	})
}

func TestFalse(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		g := ghost.New(t)

		v := true
		result := be.False(v)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, "v is true"))

		result = be.False(true)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, "true is true"))
	})

	t.Run("false", func(t *testing.T) {
		g := ghost.New(t)

		v := false
		result := be.False(v)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, "v is false"))

		result = be.False(false)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, "false is false"))
	})
}

func TestJSONEqual(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		g := ghost.New(t)

		got := `{"bar": [1, 2], "foo": "value"}`
		want := `{"foo": "value", "bar": [1, 2]}`

		result := be.JSONEqual(got, want)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, "got and want are JSON equal"))

		result = be.JSONEqual(`{"bar": [1, 2], "foo": "value"}`, `{"foo": "value", "bar": [1, 2]}`)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(
			result.Message,
			"`{\"bar\": [1, 2], \"foo\": \"value\"}` and "+
				"`{\"foo\": \"value\", \"bar\": [1, 2]}` are JSON equal",
		))
	})

	t.Run("not equal", func(t *testing.T) {
		g := ghost.New(t)

		got := `{"bar": [2, 1], "foo": "other"}`
		want := `{"foo": "value", "bar": [1, 2]}`

		result := be.JSONEqual(got, want)
		g.Should(be.False(result.Ok))
		g.Should(be.StringContaining(result.Message, "got and want are not JSON equal"))

		result = be.JSONEqual(`{"bar": [2, 1], "foo": "other"}`, `{"foo": "value", "bar": [1, 2]}`)
		g.Should(be.False(result.Ok))
		g.Should(be.StringContaining(
			result.Message,
			"`{\"bar\": [2, 1], \"foo\": \"other\"}` and "+
				"`{\"foo\": \"value\", \"bar\": [1, 2]}` are not JSON equal",
		))
	})

	t.Run("invalid json", func(t *testing.T) {
		g := ghost.New(t)

		valid := `{"foo": "value", "bar": [1, 2]}`
		invalid := `{{`
		invalid2 := `}}`

		result := be.JSONEqual(valid, invalid)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `invalid is not valid JSON
value: {{`))

		result = be.JSONEqual(invalid, valid)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `invalid is not valid JSON
value: {{`))

		result = be.JSONEqual(invalid, invalid2)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `invalid and invalid2 are not valid JSON
got:
{{

want:
}}`))

		result = be.JSONEqual(`{"bar": [1, 2], "foo": "value"}`, `{"foo": "value", "bar": [1, 2]}`)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(
			result.Message,
			"`{\"bar\": [1, 2], \"foo\": \"value\"}` and "+
				"`{\"foo\": \"value\", \"bar\": [1, 2]}` are JSON equal",
		))
	})
}

func TestMapLen(t *testing.T) {
	t.Run("equal <= 3", func(t *testing.T) {
		g := ghost.New(t)

		m := map[string]int{"a": 1, "b": 2, "c": 3}
		wantLen := 3

		result := be.MapLen(m, wantLen)
		g.Should(be.True(result.Ok))
		g.Should(be.StringContaining(result.Message, `got m length 3, want 3`))

		result = be.MapLen(map[string]int{"a": 1, "b": 2, "c": 3}, 3)
		g.Should(be.True(result.Ok))
		g.Should(be.StringContaining(
			result.Message,
			`got map[string]int{"a": 1, "b": 2, "c": 3} length 3, want 3`,
		))
	})

	t.Run("equal > 3", func(t *testing.T) {
		g := ghost.New(t)

		m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
		wantLen := 4

		result := be.MapLen(m, wantLen)
		g.Should(be.True(result.Ok))
		g.Should(be.StringContaining(result.Message, `got m length 4, want 4`))

		result = be.MapLen(map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}, 4)
		g.Should(be.True(result.Ok))
		g.Should(be.StringContaining(
			result.Message,
			`got map[string]int{"a": 1, "b": 2, "c": 3, "d": 4} length 4, want 4`,
		))
	})

	t.Run("not equal <= 3", func(t *testing.T) {
		g := ghost.New(t)

		m := map[string]int{"a": 1, "b": 2, "c": 3}
		wantLen := 2

		result := be.MapLen(m, wantLen)
		g.Should(be.False(result.Ok))
		g.Should(be.StringContaining(result.Message, `got m length 3, want 2`))

		result = be.MapLen(map[string]int{"a": 1, "b": 2, "c": 3}, 2)
		g.Should(be.False(result.Ok))
		g.Should(be.StringContaining(
			result.Message,
			`got map[string]int{"a": 1, "b": 2, "c": 3} length 3, want 2`,
		))
	})

	t.Run("not equal > 3", func(t *testing.T) {
		g := ghost.New(t)

		m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
		wantLen := 3

		result := be.MapLen(m, wantLen)
		g.Should(be.False(result.Ok))
		g.Should(be.StringContaining(result.Message, `got m length 4, want 3`))

		result = be.MapLen(map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}, 3)
		g.Should(be.False(result.Ok))
		g.Should(be.StringContaining(
			result.Message,
			`got map[string]int{"a": 1, "b": 2, "c": 3, "d": 4} length 4, want 3`,
		))
	})
}

func TestNil(t *testing.T) {
	t.Run("untyped nil", func(t *testing.T) {
		g := ghost.New(t)

		var v interface{}

		result := be.Nil(v)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, "v is nil"))

		result = be.Nil(nil)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, "nil is nil"))
	})

	t.Run("typed nil", func(t *testing.T) {
		g := ghost.New(t)

		var v *int
		var i interface{} = v

		result := be.Nil(i)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, "i is nil"))

		result = be.Nil((*int)(nil))
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, "(*int)(nil) is nil"))
	})

	t.Run("non-nil", func(t *testing.T) {
		g := ghost.New(t)

		var v int

		result := be.Nil(v)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, "v is 0, not nil"))

		result = be.Nil(-1 + 1)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, "-1 + 1 is 0, not nil"))
	})
}

func TestPanic(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		g := ghost.New(t)

		f := func() { panic(errors.New("oh no")) }

		result := be.Panic(f)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, "function f panicked with value: oh no"))

		result = be.Panic(func() { panic(errors.New("oh no")) })
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `function panicked with value: oh no
func() {
	panic(errors.New("oh no"))
}
`))
	})

	t.Run("no panic", func(t *testing.T) {
		g := ghost.New(t)

		f := func() {}

		result := be.Panic(f)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, "function f did not panic"))

		result = be.Panic(func() {})
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `function did not panic
func() {
}
`))
	})
}

func TestSliceContaining(t *testing.T) {
	t.Run("contains <= 3", func(t *testing.T) {
		g := ghost.New(t)

		slice := []int{1, 2, 3}
		elem := 2

		result := be.SliceContaining(slice, elem)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `slice contains elem
slice:   [1 2 3]
element: 2
`))

		result = be.SliceContaining([]int{1, 2, 3}, 2)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `[]int{1, 2, 3} contains 2
slice:   [1 2 3]
element: 2
`))
	})

	t.Run("contains > 3", func(t *testing.T) {
		g := ghost.New(t)

		slice := []int{1, 2, 3, 4}
		elem := 2

		result := be.SliceContaining(slice, elem)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `slice contains elem
slice:   [
	1
>	2
	3
	4
]
element: 2
`))

		result = be.SliceContaining([]int{1, 2, 3, 4}, 2)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `[]int{1, 2, 3, 4} contains 2
slice:   [
	1
>	2
	3
	4
]
element: 2
`))
	})

	t.Run("does not contain <= 3", func(t *testing.T) {
		g := ghost.New(t)

		slice := []int{1, 2, 3}
		elem := 5

		result := be.SliceContaining(slice, elem)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `slice does not contain elem
slice:   [1 2 3]
element: 5
`))

		result = be.SliceContaining([]int{1, 2, 3}, 5)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `[]int{1, 2, 3} does not contain 5
slice:   [1 2 3]
element: 5
`))
	})

	t.Run("does not contain > 3", func(t *testing.T) {
		g := ghost.New(t)

		slice := []int{1, 2, 3, 4}
		elem := 5

		result := be.SliceContaining(slice, elem)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `slice does not contain elem
slice:   [
	1
	2
	3
	4
]
element: 5
`))

		result = be.SliceContaining([]int{1, 2, 3, 4}, 5)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `[]int{1, 2, 3, 4} does not contain 5
slice:   [
	1
	2
	3
	4
]
element: 5
`))
	})
}

func TestSliceLen(t *testing.T) {
	t.Run("equal <= 3", func(t *testing.T) {
		g := ghost.New(t)

		slice := []string{"a", "b", "c"}
		wantLen := 3

		result := be.SliceLen(slice, wantLen)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `got slice length 3, want 3
slice: [a b c]
`))

		result = be.SliceLen([]string{"a", "b", "c"}, 3)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `got []string{"a", "b", "c"} length 3, want 3
slice: [a b c]
`))
	})

	t.Run("equal > 3", func(t *testing.T) {
		g := ghost.New(t)

		slice := []string{"a", "b", "c", "d"}
		wantLen := 4

		result := be.SliceLen(slice, wantLen)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `got slice length 4, want 4
slice: [
	a
	b
	c
	d
]
`))

		result = be.SliceLen([]string{"a", "b", "c", "d"}, 4)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `got []string{"a", "b", "c", "d"} length 4, want 4
slice: [
	a
	b
	c
	d
]
`))
	})

	t.Run("not equal <= 3", func(t *testing.T) {
		g := ghost.New(t)

		slice := []string{"a", "b", "c"}
		wantLen := 2

		result := be.SliceLen(slice, wantLen)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `got slice length 3, want 2
slice: [a b c]
`))

		result = be.SliceLen([]string{"a", "b", "c"}, 2)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `got []string{"a", "b", "c"} length 3, want 2
slice: [a b c]
`))
	})

	t.Run("not equal > 3", func(t *testing.T) {
		g := ghost.New(t)

		slice := []string{"a", "b", "c", "d"}
		wantLen := 3

		result := be.SliceLen(slice, wantLen)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `got slice length 4, want 3
slice: [
	a
	b
	c
	d
]
`))

		result = be.SliceLen([]string{"a", "b", "c", "d"}, 3)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `got []string{"a", "b", "c", "d"} length 4, want 3
slice: [
	a
	b
	c
	d
]
`))
	})
}

func TestStringContaining(t *testing.T) {
	t.Run("contains", func(t *testing.T) {
		g := ghost.New(t)

		outer := "foobar"
		inner := "oob"

		result := be.StringContaining(outer, inner)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `outer contains inner
str:    "foobar"
substr: "oob"
`))

		result = be.StringContaining("foobar", "oob")
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `"foobar" contains "oob"
str:    "foobar"
substr: "oob"
`))
	})

	t.Run("does not contain", func(t *testing.T) {
		g := ghost.New(t)

		outer := "foobar"
		inner := "boo"

		result := be.StringContaining(outer, inner)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `outer does not contain inner
str:    "foobar"
substr: "boo"
`))

		result = be.StringContaining("foobar", "boo")
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `"foobar" does not contain "boo"
str:    "foobar"
substr: "boo"
`))
	})

	t.Run("multiline", func(t *testing.T) {
		g := ghost.New(t)

		outer := `one
two
three
`

		result := be.StringContaining(outer, "two")
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, `outer contains "two"
str:    `+`
"""
one
two
three

"""
substr: "two"
`))
	})
}

func TestTrue(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		g := ghost.New(t)

		v := true
		result := be.True(v)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, "v is true"))

		result = be.True(true)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, "true is true"))
	})

	t.Run("false", func(t *testing.T) {
		g := ghost.New(t)

		v := false
		result := be.True(v)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, "v is false"))

		result = be.True(false)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, "false is false"))
	})
}

func TestZero(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		g := ghost.New(t)

		var v int
		result := be.Zero(v)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, "v is the zero value"))

		result = be.Zero(0)
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(result.Message, "0 is the zero value"))
	})

	t.Run("non-zero", func(t *testing.T) {
		g := ghost.New(t)

		v := 1
		result := be.Zero(v)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, "v is non-zero\nvalue: 1"))

		result = be.Zero(1)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, "1 is non-zero"))
	})
}
