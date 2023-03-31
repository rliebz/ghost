package ghost

import (
	"fmt"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/nsf/jsondiff"
)

// BeNil asserts that the given value is nil.
func BeNil(v any) Assertion {
	args := getFormattedArgs([]any{v})

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
	args := getFormattedArgs([]any{b})

	return func() Result {
		return Result{
			Success: b,
			Message: fmt.Sprintf("%v is %t", args[0], b),
		}
	}
}

// BeZero asserts that the given value equals its zero value.
func BeZero[T comparable](v T) Assertion {
	args := getFormattedArgs([]any{v})

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
	args := getFormattedArgs([]any{slice, element})

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
						sliceToString(slice, element),
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
				sliceToString(slice, element),
				element,
			),
		}
	}
}

// sliceToString pretty prints a slice, highlighting an element if it exists.
func sliceToString[T comparable](slice []T, element T) string {
	if len(slice) <= 3 {
		return fmt.Sprint(slice)
	}

	var sb strings.Builder
	sb.WriteString("[\n")
	for _, e := range slice {
		if e == element {
			sb.WriteByte('>')
		}

		sb.WriteByte('\t')
		fmt.Fprint(&sb, e)
		sb.WriteByte('\n')
	}
	sb.WriteString("]")
	return sb.String()
}

// ContainString asserts that a string contains a particular substring.
func ContainString(str, substr string) Assertion {
	args := getFormattedArgs([]any{str, substr})

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
"""`, s)
	}

	return fmt.Sprintf("%q", s)
}

// DeepEqual asserts that two elements are deeply equal.
func DeepEqual[T any](want, got T) Assertion {
	args := getFormattedArgs([]any{want, got})

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
value: %v`, args[0], args[1], want),
		}
	}
}

// Equal asserts that two elements are equal.
func Equal[T comparable](want T, got T) Assertion {
	args := getFormattedArgs([]any{want, got})

	return func() Result {
		if want == got {
			return Result{
				Success: true,
				Message: fmt.Sprintf(`%v == %v
value: %v
`, args[0], args[1], want),
			}
		}

		return Result{
			Success: false,
			Message: fmt.Sprintf(`%v != %v
diff (-want +got):
%v
`, args[0], args[1], cmp.Diff(want, got)),
		}
	}
}

// Error asserts that an error is non-nil.
func Error(err error) Assertion {
	args := getFormattedArgs([]any{err})

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

// ErrorContaining asserts that a string contains a particular substring.
func ErrorContaining(err error, msg string) Assertion {
	args := getFormattedArgs([]any{err, msg})

	return func() Result {
		if err == nil {
			return Result{
				Success: false,
				Message: "error is nil",
			}
		}

		if strings.Contains(err.Error(), msg) {
			return Result{
				Success: true,
				Message: fmt.Sprintf("%v contains message %q: %v", args[0], msg, err),
			}
		}

		return Result{
			Success: false,
			Message: fmt.Sprintf("%v does not contain message %q: %v", args[0], msg, err),
		}
	}
}

var jsonCompareOpts = jsondiff.DefaultConsoleOptions()

func JSONEqual[T ~string | ~[]byte](want, got T) Assertion {
	args := getFormattedArgs([]any{want, got})

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
