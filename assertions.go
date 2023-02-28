package ghost

import (
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/exp/slices"
)

// BeNil asserts that the given value is nil.
func BeNil(v any) Assertion {
	return func() Result {
		if v == nil {
			return Result{
				Success: true,
				Message: "value is nil",
			}
		}

		return Result{
			Success: false,
			Message: fmt.Sprintf("%q is not nil", v),
		}
	}
}

// BeTrue asserts that a value is true.
func BeTrue(b bool) Assertion {
	return func() Result {
		return Result{
			Success: b,
			Message: fmt.Sprintf("value was %t", b),
		}
	}
}

// BeZero asserts that the given value equals its zero value.
func BeZero[T comparable](v T) Assertion {
	return func() Result {
		var zero T
		if v == zero {
			return Result{
				Success: true,
				Message: "%v is the zero value",
			}
		}

		return Result{
			Success: false,
			Message: fmt.Sprintf("%v is non-zero", v),
		}
	}
}

// Contain asserts that a slice contains a particular element.
func Contain[T comparable](slice []T, element T) Assertion {
	return func() Result {
		if slices.Contains(slice, element) {
			return Result{
				Success: true,
				Message: fmt.Sprintf("%v present in %v", element, slice),
			}
		}

		return Result{
			Success: false,
			Message: fmt.Sprintf("%v not present in %v", element, slice),
		}
	}
}

// ContainString asserts that a string contains a particular substring.
func ContainString(s, substr string) Assertion {
	return func() Result {
		if strings.Contains(s, substr) {
			return Result{
				Success: true,
				Message: fmt.Sprintf("%q present in %q", substr, s),
			}
		}

		return Result{
			Success: false,
			Message: fmt.Sprintf("%q not present in %q", substr, s),
		}
	}
}

// DeepEqual asserts that two elements are deeply equal.
func DeepEqual[T any](want, got T) Assertion {
	return func() Result {
		if reflect.DeepEqual(want, got) {
			return Result{
				Success: true,
				Message: fmt.Sprintf("%v == %v", want, got),
			}
		}

		return Result{
			Success: false,
			Message: fmt.Sprintf("%v != %v", want, got),
		}
	}
}

// Equal asserts that two elements are equal.
func Equal[T comparable](want, got T) Assertion {
	return func() Result {
		if want == got {
			return Result{
				Success: true,
				Message: fmt.Sprintf("%v == %v", want, got),
			}
		}

		return Result{
			Success: false,
			Message: fmt.Sprintf("%v != %v", want, got),
		}
	}
}

// Error asserts that an error is non-nil.
func Error(err error) Assertion {
	return func() Result {
		if err == nil {
			return Result{
				Success: false,
				Message: "error was nil",
			}
		}

		return Result{
			Success: true,
			Message: fmt.Sprintf("error found with value: %q", err),
		}
	}
}

// ErrorContaining asserts that a string contains a particular substring.
func ErrorContaining(err error, msg string) Assertion {
	return func() Result {
		if err == nil {
			return Result{
				Success: false,
				Message: "error was nil",
			}
		}

		if strings.Contains(err.Error(), msg) {
			return Result{
				Success: true,
				Message: fmt.Sprintf("%q present in error %q", msg, err),
			}
		}

		return Result{
			Success: false,
			Message: fmt.Sprintf("%q not present in error %q", msg, err),
		}
	}
}

// Panic asserts that the given function panics when invoked.
func Panic(f func()) Assertion {
	return func() (result Result) {
		defer func() {
			if r := recover(); r != nil {
				result = Result{
					Success: true,
					Message: fmt.Sprintf("function panicked with value: %q", r),
				}
			}
		}()

		f()

		return Result{
			Success: false,
			Message: "function did not panic",
		}
	}
}
