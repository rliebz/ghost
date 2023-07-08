/*
Package ghost does test assertions in Go.

This is early-stage software, and some breaking changes are still expected.

# Features

  - Utilities that stay out of your way: Ghost extends standard library
    testing without trying to replace it. Set up your test in one line of code,
    and you're good to go.
  - Generics-friendly logic: All built-in assertions were designed with
    generics in mind. Shift failures left, and let your compiler tell you when
    assertions aren't set up correctly.
  - Test output that knows your code: Ghost uses AST parsing to read your
    source code and print out the most useful information for you.
  - Assertions that can be negated, extended, and reused: Easily write custom
    test assertions that are as simple to use as the built-ins. Every assertion
    can be used four different ways.

# Quick Start

Start each test by calling `g := ghost.New(t)`. Then, write your assertions:

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
*/
package ghost
