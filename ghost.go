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
func (r Runner) Should(result Result) bool {
	if h, ok := r.t.(interface{ Helper() }); ok {
		h.Helper()
	}

	if !result.Ok {
		r.t.Log(result.Message)
		r.t.Fail()
		return false
	}

	return true
}

// ShouldNot runs an assertion that should not be successful, returning true if
// the assertion was not successful.
func (r Runner) ShouldNot(result Result) bool {
	if h, ok := r.t.(interface{ Helper() }); ok {
		h.Helper()
	}

	if result.Ok {
		r.t.Log(result.Message)
		r.t.Fail()
		return false
	}

	return true
}

// Must runs an assertion that must be successful, failing the test if it is not.
func (r Runner) Must(result Result) {
	if h, ok := r.t.(interface{ Helper() }); ok {
		h.Helper()
	}

	if !r.Should(result) {
		r.t.FailNow()
	}
}

// MustNot runs an assertion that must not be successful, failing the test if it is.
func (r Runner) MustNot(result Result) {
	if h, ok := r.t.(interface{ Helper() }); ok {
		h.Helper()
	}

	if !r.ShouldNot(result) {
		r.t.FailNow()
	}
}

// NoError asserts that an error should be nil, failing the test if it is not.
func (r Runner) NoError(err error) {
	if h, ok := r.t.(interface{ Helper() }); ok {
		h.Helper()
	}

	args := getArgsFromAST([]any{err})

	if err != nil {
		r.t.Log(fmt.Sprintf("%s has error value: %s", args[0], err))
		r.t.FailNow()
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
