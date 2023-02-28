package ghost

import (
	"testing"
)

// A Runner runs test assertions.
type Runner struct {
	t *testing.T
}

// New creates a new [Runner].
func New(t *testing.T) Runner {
	return Runner{t}
}

// Should runs an assertion, returning true if the assertion was successful.
func (r Runner) Should(a Assertion) bool {
	r.t.Helper()

	result := a()

	if !result.Success {
		r.t.Error(result.Message)
		return false
	}

	return true
}

// ShouldNot runs an assertion that should not be successful, returning true if
// the assertion was not successful.
func (r Runner) ShouldNot(a Assertion) bool {
	r.t.Helper()

	result := a()

	if result.Success {
		r.t.Error(result.Message)
		return false
	}

	return true
}

// Must runs an assertion that must be successful, failing the test if it is not.
func (r Runner) Must(a Assertion) {
	r.t.Helper()

	if !r.Should(a) {
		r.t.FailNow()
	}
}

// Must not runs an assertion that must not be successful, failing the test if it is.
func (r Runner) MustNot(a Assertion) {
	r.t.Helper()

	if !r.ShouldNot(a) {
		r.t.FailNow()
	}
}

// An Assertion is any function that returns a result.
type Assertion func() Result

// A Result represents the result of an assertion.
type Result struct {
	// Success returns whether the assertion was successful.
	Success bool

	// Message returns a message describing the assertion.
	//
	// A message is required regardless of whether or not the failure was
	// successful, as results can be negated.
	Message string
}
