package be_test

import (
	"errors"
	"testing"

	"github.com/rliebz/ghost"
	"github.com/rliebz/ghost/be"
)

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
			`error err contains message msg
got:  foobar
want: oob`,
		))

		result = be.ErrorContaining(errors.New("foobar"), "oob")
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(
			result.Message,
			`error errors.New("foobar") contains message "oob"
got:  foobar
want: oob`,
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
			`error err does not contain message msg
got:  foobar
want: boo`,
		))

		result = be.ErrorContaining(errors.New("foobar"), "boo")
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(
			result.Message,
			`error errors.New("foobar") does not contain message "boo"
got:  foobar
want: boo`,
		))
	})

	t.Run("nil", func(t *testing.T) {
		g := ghost.New(t)

		var err error
		msg := "boo"

		result := be.ErrorContaining(err, msg)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `error err is nil, does not contain msg
got:  <nil>
want: boo`))

		result = be.ErrorContaining(nil, "boo")
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `error nil is nil, does not contain message
got:  <nil>
want: boo`))
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
			`error err has message msg
value: foobar`,
		))

		result = be.ErrorEqual(errors.New("foobar"), "foobar")
		g.Should(be.True(result.Ok))
		g.Should(be.Equal(
			result.Message,
			`error errors.New("foobar") has message "foobar"
value: foobar`,
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
			`error err does not have message msg
got:  foobar
want: boo`,
		))

		result = be.ErrorEqual(errors.New("foobar"), "boo")
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(
			result.Message,
			`error errors.New("foobar") does not have message "boo"
got:  foobar
want: boo`,
		))
	})

	t.Run("nil", func(t *testing.T) {
		g := ghost.New(t)

		var err error
		msg := "boo"

		result := be.ErrorEqual(err, msg)
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `error err is nil
got:  <nil>
want: boo`))

		result = be.ErrorEqual(nil, "boo")
		g.Should(be.False(result.Ok))
		g.Should(be.Equal(result.Message, `error nil is nil
got:  <nil>
want: boo`))
	})
}
