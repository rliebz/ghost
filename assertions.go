package ghost

import (
	"fmt"
	"reflect"
	"strings"
)

// BeNil asserts that the given value is nil.
func BeNil(v any) Assertion {
	fv := fmt.Sprintf("%v", v)
	if args, ok := getFormattedArgs(1); ok {
		fv = args[0]
	}

	return func() Result {
		if v == nil {
			return Result{
				Success: true,
				Message: fmt.Sprintf("%v is nil", fv),
			}
		}

		return Result{
			Success: false,
			Message: fmt.Sprintf("%v is %q, not nil", fv, v),
		}
	}
}

// BeTrue asserts that a value is true.
func BeTrue(b bool) Assertion {
	fb := fmt.Sprintf("%v", b)
	if args, ok := getFormattedArgs(1); ok {
		fb = args[0]
	}

	return func() Result {
		return Result{
			Success: b,
			Message: fmt.Sprintf("%v was %t", fb, b),
		}
	}
}

// BeZero asserts that the given value equals its zero value.
func BeZero[T comparable](v T) Assertion {
	fv := fmt.Sprintf("%v", v)
	if args, ok := getFormattedArgs(1); ok {
		fv = args[0]
	}

	return func() Result {
		var zero T
		if v == zero {
			return Result{
				Success: true,
				Message: fmt.Sprintf("%v is the zero value", fv),
			}
		}

		return Result{
			Success: false,
			Message: fmt.Sprintf("%v is non-zero", fv),
		}
	}
}

// Contain asserts that a slice contains a particular element.
func Contain[T comparable](slice []T, element T) Assertion {
	fslice, felement := fmt.Sprintf("%v", slice), fmt.Sprintf("%v", element)
	if args, ok := getFormattedArgs(1); ok {
		fslice, felement = args[0], args[1]
	}

	// TODO: Print the values of the slices / elements
	return func() Result {
		for _, x := range slice {
			if x == element {
				return Result{
					Success: true,
					Message: fmt.Sprintf("%v contains %v", fslice, felement),
				}
			}
		}

		return Result{
			Success: false,
			Message: fmt.Sprintf("%v does not contain %v", fslice, felement),
		}
	}
}

// ContainString asserts that a string contains a particular substring.
func ContainString(str, substr string) Assertion {
	fstr, fsubstr := fmt.Sprintf("%v", str), fmt.Sprintf("%v", substr)
	if args, ok := getFormattedArgs(1); ok {
		fstr, fsubstr = args[0], args[1]
	}

	// TODO: Print the values
	return func() Result {
		if strings.Contains(str, substr) {
			return Result{
				Success: true,
				Message: fmt.Sprintf("%v contains %v", fstr, fsubstr),
			}
		}

		return Result{
			Success: false,
			Message: fmt.Sprintf("%q does not contain %q", fstr, fsubstr),
		}
	}
}

// DeepEqual asserts that two elements are deeply equal.
func DeepEqual[T any](want, got T) Assertion {
	fwant, fgot := fmt.Sprintf("%v", want), fmt.Sprintf("%v", got)
	if args, ok := getFormattedArgs(1); ok {
		fwant, fgot = args[0], args[1]
	}

	// TODO: Print the values
	return func() Result {
		if reflect.DeepEqual(want, got) {
			return Result{
				Success: true,
				Message: fmt.Sprintf("%v == %v", fwant, fgot),
			}
		}

		return Result{
			Success: false,
			Message: fmt.Sprintf("%v != %v", fwant, fgot),
		}
	}
}

// Equal asserts that two elements are equal.
func Equal[T comparable](want T, got T) Assertion {
	fwant, fgot := fmt.Sprintf("%v", want), fmt.Sprintf("%v", got)
	if args, ok := getFormattedArgs(1); ok {
		fwant, fgot = args[0], args[1]
	}

	// TODO: Print the values
	return func() Result {
		if want == got {
			return Result{
				Success: true,
				Message: fmt.Sprintf("%v == %v", fwant, fgot),
			}
		}

		return Result{
			Success: false,
			Message: fmt.Sprintf("%v != %v", fwant, fgot),
		}
	}
}

// Error asserts that an error is non-nil.
func Error(err error) Assertion {
	ferr := fmt.Sprintf("%v", err)
	if args, ok := getFormattedArgs(1); ok {
		ferr = args[0]
	}

	return func() Result {
		if err == nil {
			return Result{
				Success: false,
				Message: fmt.Sprintf("%s was nil", ferr),
			}
		}

		return Result{
			Success: true,
			Message: fmt.Sprintf("%s contains error: %s", ferr, err),
		}
	}
}

// ErrorContaining asserts that a string contains a particular substring.
func ErrorContaining(err error, msg string) Assertion {
	ferr := fmt.Sprintf("%v", err)
	if args, ok := getFormattedArgs(1); ok {
		ferr = args[0]
	}

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
				Message: fmt.Sprintf("%v contains message %q: %v", ferr, msg, err),
			}
		}

		return Result{
			Success: false,
			Message: fmt.Sprintf("%v does not contain message %q: %v", ferr, msg, err),
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
