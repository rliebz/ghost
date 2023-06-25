package ghost

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/nsf/jsondiff"

	"github.com/rliebz/ghost/internal/constraints"
)

// BeInDelta asserts that a value is within a delta of another.
func BeInDelta[T constraints.Integer | constraints.Float](want, got, delta T) Result {
	args := getArgsFromAST([]any{want, got, delta})

	diff := want - got
	if diff < 0 {
		diff = 0 - diff
	}

	wantStr := args[0]
	if _, err := strconv.ParseFloat(wantStr, 64); err != nil {
		wantStr = fmt.Sprintf("%s (%v)", wantStr, want)
	}

	gotStr := args[1]
	if _, err := strconv.ParseFloat(gotStr, 64); err != nil {
		gotStr = fmt.Sprintf("%s (%v)", gotStr, got)
	}

	if diff <= delta {
		return Result{
			Ok:      true,
			Message: fmt.Sprintf("delta %v between %s and %s is within %v", diff, wantStr, gotStr, delta),
		}
	}

	return Result{
		Ok:      false,
		Message: fmt.Sprintf("delta %v between %s and %s is not within %v", diff, wantStr, gotStr, delta),
	}
}

// BeNil asserts that the given value is nil.
func BeNil(v any) Result {
	args := getArgsFromAST([]any{v})

	if v == nil {
		return Result{
			Ok:      true,
			Message: fmt.Sprintf("%v is nil", args[0]),
		}
	}

	return Result{
		Ok:      false,
		Message: fmt.Sprintf("%v is %v, not nil", args[0], v),
	}
}

// BeTrue asserts that a value is true.
func BeTrue(b bool) Result {
	args := getArgsFromAST([]any{b})

	return Result{
		Ok:      b,
		Message: fmt.Sprintf("%v is %t", args[0], b),
	}
}

// BeZero asserts that the given value equals its zero value.
func BeZero[T comparable](v T) Result {
	args := getArgsFromAST([]any{v})

	var zero T
	if v == zero {
		return Result{
			Ok:      true,
			Message: fmt.Sprintf("%v is the zero value", args[0]),
		}
	}

	if args[0] != fmt.Sprintf("%v", v) {
		return Result{
			Ok:      false,
			Message: fmt.Sprintf("%v is non-zero\nvalue: %v", args[0], v),
		}
	}

	return Result{
		Ok:      false,
		Message: fmt.Sprintf("%v is non-zero", args[0]),
	}
}

// Contain asserts that a slice contains a particular element.
func Contain[T comparable](slice []T, element T) Result {
	args := getArgsFromAST([]any{slice, element})

	for _, x := range slice {
		if x == element {
			return Result{
				Ok: true,
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
		Ok: false,
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

// ContainString asserts that a string contains a particular substring.
func ContainString(str, substr string) Result {
	args := getArgsFromAST([]any{str, substr})

	if strings.Contains(str, substr) {
		return Result{
			Ok: true,
			Message: fmt.Sprintf(`%v contains %v
str:    %s
substr: %s
`, args[0], args[1], quoteString(str), quoteString(substr)),
		}
	}

	return Result{
		Ok: false,
		Message: fmt.Sprintf(`%v does not contain %v
str:    %s
substr: %s
`, args[0], args[1], quoteString(str), quoteString(substr)),
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
func DeepEqual[T any](want, got T) Result {
	args := getArgsFromAST([]any{want, got})

	if diff := cmp.Diff(want, got); diff != "" {
		return Result{
			Ok: false,
			Message: fmt.Sprintf(`%v != %v
diff (-want +got):
%v
`, args[0], args[1], diff),
		}
	}

	return Result{
		Ok: true,
		Message: fmt.Sprintf(`%v == %v
value: %v
`, args[0], args[1], want),
	}
}

// Equal asserts that two elements are equal.
func Equal[T comparable](want T, got T) Result {
	args := getArgsFromAST([]any{want, got})

	if want == got {
		return Result{
			Ok: true,
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
			Ok: false,
			Message: fmt.Sprintf(`%v != %v
diff (-want +got):
%v
`, args[0], args[1], cmp.Diff(want, got)),
		}
	case reflect.String:
		return Result{
			Ok: false,
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
		Ok: false,
		Message: fmt.Sprintf(`%v != %v
want: %v
got:  %v
`, args[0], args[1], want, got),
	}
}

// Error asserts that an error is non-nil.
func Error(err error) Result {
	args := getArgsFromAST([]any{err})

	if err == nil {
		return Result{
			Ok:      false,
			Message: fmt.Sprintf("%s is nil", args[0]),
		}
	}

	return Result{
		Ok:      true,
		Message: fmt.Sprintf("%s has error value: %s", args[0], err),
	}
}

// ErrorContaining asserts that an error string contains a particular substring.
func ErrorContaining(msg string, err error) Result {
	args := getArgsFromAST([]any{msg, err})

	switch {
	case err == nil && args[0] == fmt.Sprintf("%q", msg):
		return Result{
			Ok:      false,
			Message: fmt.Sprintf(`%v is nil; missing error message: %v`, args[1], msg),
		}
	case err == nil:
		return Result{
			Ok:      false,
			Message: fmt.Sprintf(`%v is nil; missing error message %v: %v`, args[1], args[0], msg),
		}
	case strings.Contains(err.Error(), msg):
		return Result{
			Ok:      true,
			Message: fmt.Sprintf("%v contains error message %q: %v", args[1], msg, err),
		}
	default:
		return Result{
			Ok:      false,
			Message: fmt.Sprintf("%v does not contain error message %q: %v", args[1], msg, err),
		}
	}
}

// ErrorEqual asserts that an error string equals a particular message.
func ErrorEqual(msg string, err error) Result {
	args := getArgsFromAST([]any{msg, err})

	if err == nil {
		return Result{
			Ok:      false,
			Message: fmt.Sprintf(`%v is nil; want message: %v`, args[1], msg),
		}
	}

	if err.Error() == msg {
		return Result{
			Ok:      true,
			Message: fmt.Sprintf("%v equals error message %q: %v", args[1], msg, err),
		}
	}

	return Result{
		Ok:      false,
		Message: fmt.Sprintf("%v does not equal error message %q: %v", args[1], msg, err),
	}
}

var jsonCompareOpts = jsondiff.DefaultConsoleOptions()

// JSONEqual asserts that two sets of JSON-encoded data are equivalent.
func JSONEqual[T ~string | ~[]byte](want, got T) Result {
	args := getArgsFromAST([]any{want, got})

	diff, msg := jsondiff.Compare([]byte(want), []byte(got), &jsonCompareOpts)

	switch diff {
	case jsondiff.FullMatch:
		return Result{
			Ok:      true,
			Message: fmt.Sprintf("%v and %v are JSON equal", args[0], args[1]),
		}
	case jsondiff.FirstArgIsInvalidJson:
		return Result{
			Ok: false,
			Message: fmt.Sprintf(`%v is not valid JSON
value: %s`, args[0], want),
		}
	case jsondiff.SecondArgIsInvalidJson:
		return Result{
			Ok: false,
			Message: fmt.Sprintf(`%v is not valid JSON
value: %s`, args[1], got),
		}
	case jsondiff.BothArgsAreInvalidJson:
		return Result{
			Ok: false,
			Message: fmt.Sprintf(`%v and %v are not valid JSON
want:
%s

got:
%s`, args[0], args[1], want, got),
		}
	}

	return Result{
		Ok:      false,
		Message: msg,
	}
}

// Len asserts that the length of a slice is a particular size.
func Len[T any](want int, got []T) Result {
	args := getArgsFromAST([]any{want, got})

	return Result{
		Ok: want == len(got),
		Message: fmt.Sprintf(`want %v length %d, got %d
slice: %v
`, args[1], want, len(got), sliceToString(got)),
	}
}

// Panic asserts that the given function panics when invoked.
func Panic(f func()) (result Result) {
	args := getArgsFromAST([]any{f})

	defer func() {
		if r := recover(); r != nil {
			if strings.Contains(args[0], "\n") {
				result = Result{
					Ok: true,
					Message: fmt.Sprintf(`function panicked with value: %v
%v
`, r, args[0]),
				}
			} else {
				result = Result{
					Ok:      true,
					Message: fmt.Sprintf(`function %v panicked with value: %v`, args[0], r),
				}
			}
		}
	}()

	f()

	if strings.Contains(args[0], "\n") {
		return Result{
			Ok: false,
			Message: fmt.Sprintf(`function did not panic
%v
`, args[0]),
		}
	}

	return Result{
		Ok:      false,
		Message: fmt.Sprintf("function %v did not panic", args[0]),
	}
}
