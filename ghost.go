package ghost

import (
	"testing"
)

// A Tester runs test assertions.
type Tester struct {
	t *testing.T
}

// New creates a new [Tester].
func New(t *testing.T) Tester {
	return Tester{t}
}

// Should runs an assertion, returning true if the assertion was successful.
func (t Tester) Should(a Assertion) bool {
	t.t.Helper()

	result := a()

	if !result.Success {
		t.t.Error(result.Message)
		return false
	}

	return true
}

// ShouldNot runs an assertion that should not be successful, returning true if
// the assertion was not successful.
func (t Tester) ShouldNot(a Assertion) bool {
	t.t.Helper()

	result := a()

	if result.Success {
		t.t.Error(result.Message)
		return false
	}

	return true
}

// Must runs an assertion that must be successful, failing the test if it is not.
func (t Tester) Must(a Assertion) {
	t.t.Helper()

	if !t.Should(a) {
		t.t.FailNow()
	}
}

// Must not runs an assertion that must not be successful, failing the test if it is.
func (t Tester) MustNot(a Assertion) {
	t.t.Helper()

	if !t.ShouldNot(a) {
		t.t.FailNow()
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
