package ghost

import "fmt"

// T is the subset of [*testing.T] used in assertions.
//
// The Helper() method will be called if available.
type T interface {
	Log(args ...any)
	Fail()
	FailNow()
}

// A Runner runs test assertions.
type Runner struct {
	t T
}

// New creates a new [Runner].
func New(t T) Runner {
	return Runner{t}
}

// Should runs an assertion, returning true if the assertion was successful.
func (r Runner) Should(a Assertion) bool {
	if h, ok := r.t.(interface{ Helper() }); ok {
		h.Helper()
	}

	result := a()

	if !result.Success {
		r.t.Log(result.Message)
		r.t.Fail()
		return false
	}

	return true
}

// ShouldNot runs an assertion that should not be successful, returning true if
// the assertion was not successful.
func (r Runner) ShouldNot(a Assertion) bool {
	if h, ok := r.t.(interface{ Helper() }); ok {
		h.Helper()
	}

	result := a()

	if result.Success {
		r.t.Log(result.Message)
		r.t.Fail()
		return false
	}

	return true
}

// Must runs an assertion that must be successful, failing the test if it is not.
func (r Runner) Must(a Assertion) {
	if h, ok := r.t.(interface{ Helper() }); ok {
		h.Helper()
	}

	if !r.Should(a) {
		r.t.FailNow()
	}
}

// MustNot runs an assertion that must not be successful, failing the test if it is.
func (r Runner) MustNot(a Assertion) {
	if h, ok := r.t.(interface{ Helper() }); ok {
		h.Helper()
	}

	if !r.ShouldNot(a) {
		r.t.FailNow()
	}
}

// NoErr asserts that an error should be nil, failing the test if it is not.
func (r Runner) NoErr(err error) {
	if h, ok := r.t.(interface{ Helper() }); ok {
		h.Helper()
	}

	args := getFormattedArgs([]any{err})

	if err != nil {
		r.t.Log(fmt.Sprintf("%s has error value: %s", args[0], err))
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
