package ghost

import (
	"fmt"

	"github.com/rliebz/ghost/ghostlib"
)

// T is the subset of [*testing.T] used in assertions.
//
// The Helper() method will be called if available.
type T interface {
	Log(args ...any)
	Fail()
	FailNow()
}

// Ghost runs test assertions.
type Ghost struct {
	t T
}

// New creates a new [Ghost].
func New(t T) Ghost {
	return Ghost{t}
}

// Should runs an assertion, returning true if the assertion was successful.
func (g Ghost) Should(result Result) bool {
	if h, ok := g.t.(interface{ Helper() }); ok {
		h.Helper()
	}

	if !result.Ok {
		g.t.Log(result.Message)
		g.t.Fail()
		return false
	}

	return true
}

// ShouldNot runs an assertion that should not be successful, returning true if
// the assertion was not successful.
func (g Ghost) ShouldNot(result Result) bool {
	if h, ok := g.t.(interface{ Helper() }); ok {
		h.Helper()
	}

	if result.Ok {
		g.t.Log(result.Message)
		g.t.Fail()
		return false
	}

	return true
}

// Must runs an assertion that must be successful, failing the test if it is not.
func (g Ghost) Must(result Result) {
	if h, ok := g.t.(interface{ Helper() }); ok {
		h.Helper()
	}

	if !g.Should(result) {
		g.t.FailNow()
	}
}

// MustNot runs an assertion that must not be successful, failing the test if it is.
func (g Ghost) MustNot(result Result) {
	if h, ok := g.t.(interface{ Helper() }); ok {
		h.Helper()
	}

	if !g.ShouldNot(result) {
		g.t.FailNow()
	}
}

// NoError asserts that an error should be nil, failing the test if it is not.
func (g Ghost) NoError(err error) {
	if h, ok := g.t.(interface{ Helper() }); ok {
		h.Helper()
	}

	args := ghostlib.ArgsFromAST(err)

	if err != nil {
		g.t.Log(fmt.Sprintf("%s has error value: %s", args[0], err))
		g.t.FailNow()
	}
}

// An Result represents the result of an assertion.
type Result struct {
	// Ok returns whether the assertion was successful.
	Ok bool

	// Message returns a message describing the assertion.
	//
	// A message is required regardless of whether or not the failure was
	// successful.
	Message string
}
