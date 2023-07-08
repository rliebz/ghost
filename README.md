# Ghost

[![Go Reference](https://pkg.go.dev/badge/github.com/rliebz/ghost.svg)][godoc]

Ghost does test assertions in Go.

This is early-stage software, and some breaking changes are still expected.

## Features

- **Utilities that stay out of your way**: Ghost extends standard library
  testing without trying to replace it. Set up your test in one line of code,
  and you're good to go.
- **Generics-friendly logic**: All built-in assertions were designed with
  generics in mind. Shift failures left, and let your compiler tell you when
  assertions aren't set up correctly.
- **Test output that knows your code**: Ghost uses AST parsing to read your
  source code and print out the most useful information for you.
- **Assertions that can be negated, extended, and reused**: Easily write custom
  test assertions that are as simple to use as the built-ins. Every assertion
  can be used four different ways.

## Quick Start

Start each test by calling `g := ghost.New(t)`. Then, write your assertions:

```go
func TestMyFunc(t *testing.T) {
  g := ghost.New(t)

  got, err := MyFunc()
  g.NoError(err)

  g.MustNot(be.Nil(got))
  g.Should(be.Equal("my value", got.SomeString))
  g.Should(be.SliceLen(3, got.SomeSlice))
}

func TestMyFunc_error(t *testing.T) {
  g := ghost.New(t)

  got, err := MyFunc()

  g.Should(be.Zero(got))
  g.Should(be.ErrorEqual("an error occurred", err))
}
```

## Usage

### Checks

Ghost comes with four main checks: `Should`, `ShouldNot`, `Must`, and `MustNot`.

`Should` and `ShouldNot` check whether an assertion has succeeded, failing the
test otherwise. Like `t.Error`, the test is allowed to proceed if the assertion
fails:

```go
g.Should(be.Equal(want, got))
```

Both functions also return a boolean indicating whether the check was
successful, allowing you to safely chain assertion logic:

```go
if g.Should(be.Len(1, mySlice)) {
  g.Should(be.Equal("foo", mySlice[0]))
}
```

`Must` and `MustNot` work similarly, but end test execution if the assertion
does not pass, analogous to `t.Fatal`:

```go
g.Must(be.True(ok))
g.MustNot(be.Nil(val))
```

For convenience, a `NoError` check is also available, which fails and ends test
execution for non-nil errors:

```go
g.NoError(err)

// Equivalent to:
g.MustNot(be.Error(err))
```

### Assertions

An assertion is any function that returns a `ghost.Result`.

#### Standard Assertions

A set of standard assertions are available in [github.com/rliebz/ghost/be][godoc/be].

These cover common use cases, such as simple and deep equality, slice/map/string
operations, error and panic handling, and JSON equality.

```go
g.Should(be.True(true))
g.ShouldNot(be.False(true))

g.Should(be.Equal(1+1, 2))
g.Should(be.DeepEqual([]string{"a", "b"}, []string{"a", "b"}))

g.Should(be.InSlice(2, []int{1, 2, 3}))
g.Should(be.InString("foo", "foobar"))

g.Should(be.Panic(func() { panic("oh no") }))

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
```

For the full list available, see [the documentation][godoc/be].

#### Custom Assertions

Custom assertions are easy to write and easy to use.

A very basic custom assertion might look like this:

```go
func BeThirteen(i int) ghost.Result {
	return ghost.Result{
		Ok:      i == 13,
		Message: fmt.Sprintf("value is %d", i),
	}
}
```

Custom assertions can be used with the built-in checks:

```go
i := 13
g.Should(BeThirteen(i))
```

One of the key features that makes test output readable is understanding the
AST to be able to print better failure messages. Use [ghostlib][godoc/ghostlib]
to pretty print the AST values of assertion arguments:

```go
func BeThirteen(i int) ghost.Result {
	args := ghostlib.ArgsFromAST(i)

	return ghost.Result{
		Ok:      i == 13,
		Message: fmt.Sprintf("%v is %d", args[0], i),
	}
}
```

And instantly get helpful, descriptive error messages:

```go
g.Should(BeThirteen(myInt)) // "myInt is 0"
g.Should(BeThirteen(5 + 6)) // "5 + 6 is 11"
```

## Philosophy

### Both "Hard" and "Soft" Assertions Should Be Easy

Some testing libraries lock you into stopping test execution on assertion
failure. Ghost makes it easy to switch between both, and doesn't make you
change the way you set tests up based on which of the two you use:

```go
g := ghost.New(t)             // universal test setup

// ...

g.Should(be.Equal(13, myInt)) // soft assertion
g.Must(be.True(ok))           // hard assertion
```

### Arguments Should Be Predictable

Arguments to assertions should go in a predictable order. By convention:

1. "Want" comes before "got".
2. "Needle" comes before "haystack".
3. All other arguments come last.

### Ghost Does Assertions

Go's `testing` package is fantastic; Ghost doesn't try to do anything that the
standard library already does.

Test suites, mocking, logging, and non-assertion failures are all out of scope.

[godoc]: https://pkg.go.dev/github.com/rliebz/ghost
[godoc/be]: https://pkg.go.dev/github.com/rliebz/ghost/be
[godoc/ghostlib]: https://pkg.go.dev/github.com/rliebz/ghost/ghostlib
