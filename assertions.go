package ghost

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/nsf/jsondiff"
)

// BeNil asserts that the given value is nil.
func BeNil(v any) Assertion {
	args := getArgsFromAST([]any{v})

	return func() Result {
		if v == nil {
			return Result{
				Success: true,
				Message: fmt.Sprintf("%v is nil", args[0]),
			}
		}

		return Result{
			Success: false,
			Message: fmt.Sprintf("%v is %v, not nil", args[0], v),
		}
	}
}

// BeTrue asserts that a value is true.
func BeTrue(b bool) Assertion {
	args := getArgsFromAST([]any{b})

	return func() Result {
		return Result{
			Success: b,
			Message: fmt.Sprintf("%v is %t", args[0], b),
		}
	}
}

// BeZero asserts that the given value equals its zero value.
func BeZero[T comparable](v T) Assertion {
	args := getArgsFromAST([]any{v})

	return func() Result {
		var zero T
		if v == zero {
			return Result{
				Success: true,
				Message: fmt.Sprintf("%v is the zero value", args[0]),
			}
		}

		if args[0] != fmt.Sprintf("%v", v) {
			return Result{
				Success: false,
				Message: fmt.Sprintf("%v is non-zero\nvalue: %v", args[0], v),
			}
		}

		return Result{
			Success: false,
			Message: fmt.Sprintf("%v is non-zero", args[0]),
		}
	}
}

// Contain asserts that a slice contains a particular element.
func Contain[T comparable](slice []T, element T) Assertion {
	args := getArgsFromAST([]any{slice, element})

	return func() Result {
		for _, x := range slice {
			if x == element {
				return Result{
					Success: true,
					Message: fmt.Sprintf(`%v contains %v
slice:   %v
element: %v
`,
						args[0],
						args[1],
						sliceElementToString(slice, element),
						element,
					),
				}
			}
		}

		return Result{
			Success: false,
			Message: fmt.Sprintf(`%v does not contain %v
slice:   %v
element: %v
`,
				args[0],
				args[1],
				sliceElementToString(slice, element),
				element,
			),
		}
	}
}

// ContainString asserts that a string contains a particular substring.
func ContainString(str, substr string) Assertion {
	args := getArgsFromAST([]any{str, substr})

	return func() Result {
		if strings.Contains(str, substr) {
			return Result{
				Success: true,
				Message: fmt.Sprintf(`%v contains %v
str:    %s
substr: %s
`, args[0], args[1], quoteString(str), quoteString(substr)),
			}
		}

		return Result{
			Success: false,
			Message: fmt.Sprintf(`%v does not contain %v
str:    %s
substr: %s
`, args[0], args[1], quoteString(str), quoteString(substr)),
		}
	}
}

// quoteString prints a string as a single quoted line, or multiline block.
func quoteString(s string) string {
	if strings.ContainsAny(s, "\n\r") {
		return fmt.Sprintf(`
"""
%s
"""
`, s)
	}

	return fmt.Sprintf("%q", s)
}

// DeepEqual asserts that two elements are deeply equal.
func DeepEqual[T any](want, got T) Assertion {
	args := getArgsFromAST([]any{want, got})

	return func() Result {
		if diff := cmp.Diff(want, got); diff != "" {
			return Result{
				Success: false,
				Message: fmt.Sprintf(`%v != %v
diff (-want +got):
%v
`, args[0], args[1], diff),
			}
		}

		return Result{
			Success: true,
			Message: fmt.Sprintf(`%v == %v
value: %v
`, args[0], args[1], want),
		}
	}
}

// Equal asserts that two elements are equal.
func Equal[T comparable](want T, got T) Assertion {
	args := getArgsFromAST([]any{want, got})

	return func() Result {
		if want == got {
			return Result{
				Success: true,
				Message: fmt.Sprintf(`%v == %v
value: %v
`, args[0], args[1], want),
			}
		}

		switch v := reflect.ValueOf(want); v.Kind() {
		// These are types cmp tends to do particularly well
		case
			reflect.Array,
			reflect.Interface,
			reflect.Map,
			reflect.Pointer,
			reflect.Slice,
			reflect.Struct:

			return Result{
				Success: false,
				Message: fmt.Sprintf(`%v != %v
diff (-want +got):
%v
`, args[0], args[1], cmp.Diff(want, got)),
			}
		case reflect.String:
			return Result{
				Success: false,
				Message: fmt.Sprintf(`%v != %v
want: %v
got:  %v
`,
					args[0],
					args[1],
					quoteString(interface{}(want).(string)),
					quoteString(interface{}(got).(string)),
				),
			}
		}

		return Result{
			Success: false,
			Message: fmt.Sprintf(`%v != %v
want: %v
got:  %v
`, args[0], args[1], want, got),
		}
	}
}

// Error asserts that an error is non-nil.
func Error(err error) Assertion {
	args := getArgsFromAST([]any{err})

	return func() Result {
		if err == nil {
			return Result{
				Success: false,
				Message: fmt.Sprintf("%s is nil", args[0]),
			}
		}

		return Result{
			Success: true,
			Message: fmt.Sprintf("%s has error value: %s", args[0], err),
		}
	}
}

// ErrorContaining asserts that an error string contains a particular substring.
func ErrorContaining(msg string, err error) Assertion {
	args := getArgsFromAST([]any{msg, err})

	return func() Result {
		switch {
		case err == nil && args[0] == fmt.Sprintf("%q", msg):
			return Result{
				Success: false,
				Message: fmt.Sprintf(`%v is nil; missing error message: %v`, args[1], msg),
			}
		case err == nil:
			return Result{
				Success: false,
				Message: fmt.Sprintf(`%v is nil; missing error message %v: %v`, args[1], args[0], msg),
			}
		case strings.Contains(err.Error(), msg):
			return Result{
				Success: true,
				Message: fmt.Sprintf("%v contains error message %q: %v", args[1], msg, err),
			}
		default:
			return Result{
				Success: false,
				Message: fmt.Sprintf("%v does not contain error message %q: %v", args[1], msg, err),
			}
		}
	}
}

// ErrorEqual asserts that an error string equals a particular message.
func ErrorEqual(msg string, err error) Assertion {
	args := getArgsFromAST([]any{msg, err})

	return func() Result {
		if err == nil {
			return Result{
				Success: false,
				Message: fmt.Sprintf(`error is nil
want message: %v
`, msg),
			}
		}

		if err.Error() == msg {
			return Result{
				Success: true,
				Message: fmt.Sprintf("%v equals error message %q: %v", args[1], msg, err),
			}
		}

		return Result{
			Success: false,
			Message: fmt.Sprintf("%v does not equal error message %q: %v", args[1], msg, err),
		}
	}
}

var jsonCompareOpts = jsondiff.DefaultConsoleOptions()

func JSONEqual[T ~string | ~[]byte](want, got T) Assertion {
	args := getArgsFromAST([]any{want, got})

	return func() Result {
		diff, msg := jsondiff.Compare([]byte(want), []byte(got), &jsonCompareOpts)

		switch diff {
		case jsondiff.FullMatch:
			return Result{
				Success: true,
				Message: fmt.Sprintf("%v and %v are JSON equal", args[0], args[1]),
			}
		case jsondiff.FirstArgIsInvalidJson:
			return Result{
				Success: false,
				Message: fmt.Sprintf(`%v is not valid JSON
value: %s`, args[0], want),
			}
		case jsondiff.SecondArgIsInvalidJson:
			return Result{
				Success: false,
				Message: fmt.Sprintf(`%v is not valid JSON
value: %s`, args[1], got),
			}
		case jsondiff.BothArgsAreInvalidJson:
			return Result{
				Success: false,
				Message: fmt.Sprintf(`%v and %v are not valid JSON
want:
%s

got:
%s`, args[0], args[1], want, got),
			}
		}

		return Result{
			Success: false,
			Message: msg,
		}
	}
}

func Len[T any](want int, got []T) Assertion {
	args := getArgsFromAST([]any{want, got})

	return func() Result {
		return Result{
			Success: want == len(got),
			Message: fmt.Sprintf(`want %v length %d, got %d
slice: %v
`, args[1], want, len(got), sliceToString(got)),
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
