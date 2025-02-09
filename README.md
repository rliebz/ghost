# Ghost

[![Go Reference](https://pkg.go.dev/badge/github.com/rliebz/ghost.svg)][godoc]
[![Test Workflow](https://github.com/rliebz/ghost/actions/workflows/test.yml/badge.svg)](https://github.com/rliebz/ghost/actions?query=workflow%3ATest+branch%3Amain++)
[![Go Report Card](https://goreportcard.com/badge/github.com/rliebz/ghost)](https://goreportcard.com/report/github.com/rliebz/ghost)

Flexible, extensible, and beautiful test assertions in Go that stay out of
your way.

## Features

- **Assertions that fit in your existing tests**: Ghost extends standard
  library testing without trying to replace it. Set up your test in one line of
  code, and you're good to go.
- **Type-safe and generics-friendly logic**: All built-in assertions were
  designed with generics in mind. Shift failures left, and let your compiler
  tell you when assertions aren't set up correctly.
- **Test output that knows your code**: Ghost uses AST parsing to read your
  source code and print out the most useful information for you.
- **Assertions that can be composed, extended, and reused**: Easily write
  custom test assertions that are as simple to use as the built-ins. Every
  assertion can transformed to express whatever you need to.

## Quick Start

Start each test by calling `g := ghost.New(t)`. Then, write your assertions:

```go
func TestMyFunc(t *testing.T) {
  g := ghost.New(t)

  got, err := MyFunc()
  g.NoError(err)

  g.MustNot(be.Nil(got))
  g.Should(be.Equal(got.SomeString, "my value"))
  g.Should(be.SliceLen(got.SomeSlice, 3))
}

func TestMyFunc_error(t *testing.T) {
  g := ghost.New(t)

  got, err := MyFunc()

  g.Should(be.Zero(got))
  g.Should(be.ErrorEqual(err, "an error occurred"))
}
```

## Usage

### Checks

Ghost comes with four main checks: `Should`, `ShouldNot`, `Must`, and `MustNot`.

`Should` and `ShouldNot` check whether an assertion has succeeded, failing the
test otherwise. Like `t.Error`, the test is allowed to proceed if the assertion
fails:

```go
g.Should(be.Equal(got, want))
g.ShouldNot(be.Nil(val))
```

Both functions also return a boolean indicating whether the check was
successful, allowing you to safely chain assertion logic:

```go
if g.Should(be.SliceLen(mySlice, 1)) {
  g.Should(be.Equal(mySlice[0], "foo"))
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

g.Should(be.SliceContaining([]int{1, 2, 3}, 2))
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
```

For the full list available, see [the documentation][godoc/be].

#### Assertion Composers

Ghost allows assertions to be composed into powerful expressions.

One composer is `be.Eventually`, which retries an assertion over time until
it either succeeds or times out:

```go
g.Should(be.Eventually(func() ghost.Result {
  return be.True(val.IsSettled(ctx))
}, 3*time.Second, 100*time.Millisecond))
```

Another composer is `be.Not`, which negates the result of an assertion:

```go
g.Should(be.Not(be.True(ok)))
g.Must(be.Not(be.Nil(val)))
```

While `be.Not` in a simple assertion would simply be a more verbose version of
of `ShouldNot` or `MustNot`, the real benefit becomes obvious when you combine
composers together:

```go
g.Should(be.Eventually(func() ghost.Result {
  return be.Not(be.Equal(a, b))
}, 3*time.Second, 100*time.Millisecond))
```

For details on other composers such as `be.Any` or `be.All`, see the [godoc][].

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

#### Handling Panics

If you expect your code to panic, it is better to assert that the value passed
to `panic` has the properties you expect, rather than to make an assumption
that the panic you encountered is the panic you were expecting. Ghost can be
combined with `defer`/`recover` to access the full expressiveness of test
assertions:

```go
defer func() {
	var err error
	g.Must(be.AssignedAs(recover(), &err))
	g.Should(be.ErrorEqual(err, "a specific error occurred"))
}()

doStuff()
```

## Philosophy

### Ghost Does Assertions

Go's `testing` package is fantastic; Ghost doesn't try to do anything that the
standard library already does.

Test suites, mocking, logging, and non-assertion failures are all out of scope.

### Both "Hard" and "Soft" Assertions Should Be Easy

Some testing libraries lock you into stopping test execution on assertion
failure. Ghost makes it easy to switch between both, and doesn't make you
change the way you set tests up based on which of the two you use:

```go
g := ghost.New(t)             // universal test setup

// ...

g.Should(be.Equal(myInt, 13)) // soft assertion
g.Must(be.True(ok))           // hard assertion
```

### Arguments Should Be Predictable

Arguments to assertions should go in an intuitive, predictable order. By
convention:

1. "Got" comes before "want".
2. "Haystack" comes before "needle".
3. All other arguments come last.

[godoc]: https://pkg.go.dev/github.com/rliebz/ghost
[godoc/be]: https://pkg.go.dev/github.com/rliebz/ghost/be
[godoc/ghostlib]: https://pkg.go.dev/github.com/rliebz/ghost/ghostlib
