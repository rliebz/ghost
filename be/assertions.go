package be

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/nsf/jsondiff"

	"github.com/rliebz/ghost"
	"github.com/rliebz/ghost/ghostlib"
	"github.com/rliebz/ghost/internal/constraints"
)

// Close asserts that a value is within a delta of another.
func Close[T constraints.Integer | constraints.Float](want, got, delta T) ghost.Result {
	args := ghostlib.ArgsFromAST(want, got, delta)
	argWant, argGot := args[0], args[1]

	diff := want - got
	if diff < 0 {
		diff = 0 - diff
	}

	if _, err := strconv.ParseFloat(argWant, 64); err != nil {
		argWant = fmt.Sprintf("%s (%v)", argWant, want)
	}

	if _, err := strconv.ParseFloat(argGot, 64); err != nil {
		argGot = fmt.Sprintf("%s (%v)", argGot, got)
	}

	if diff <= delta {
		return ghost.Result{
			Ok: true,
			Message: fmt.Sprintf(
				"delta %v between %s and %s is within %v",
				diff, argWant, argGot, delta,
			),
		}
	}

	return ghost.Result{
		Ok: false,
		Message: fmt.Sprintf(
			"delta %v between %s and %s is not within %v",
			diff, argWant, argGot, delta,
		),
	}
}

// DeepEqual asserts that two elements are deeply equal.
func DeepEqual[T any](want, got T) ghost.Result {
	args := ghostlib.ArgsFromAST(want, got)
	argWant, argGot := args[0], args[1]

	if diff := cmp.Diff(
		want, got,
		cmp.Exporter(func(reflect.Type) bool { return true }),
	); diff != "" {
		return ghost.Result{
			Ok: false,
			Message: fmt.Sprintf(`%v != %v
diff (-want +got):
%v
`, argWant, argGot, diff),
		}
	}

	return ghost.Result{
		Ok: true,
		Message: fmt.Sprintf(`%v == %v
value: %v
`, argWant, argGot, want),
	}
}

// Equal asserts that two elements are equal.
func Equal[T comparable](want T, got T) ghost.Result {
	args := ghostlib.ArgsFromAST(want, got)
	argWant, argGot := args[0], args[1]

	if want == got {
		switch fmt.Sprint(want) {
		case argWant, argGot:
			return ghost.Result{
				Ok:      true,
				Message: fmt.Sprintf(`%v == %v`, argWant, argGot),
			}
		default:
			return ghost.Result{
				Ok: true,
				Message: fmt.Sprintf(`%v == %v
value: %v
`, argWant, argGot, want),
			}
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

		return ghost.Result{
			Ok: false,
			Message: fmt.Sprintf(`%v != %v
diff (-want +got):
%v
`, argWant, argGot, cmp.Diff(want, got)),
		}
	case reflect.String:
		return ghost.Result{
			Ok: false,
			Message: fmt.Sprintf(`%v != %v
want: %v
got:  %v
`,
				argWant,
				argGot,
				quoteString(interface{}(want).(string)),
				quoteString(interface{}(got).(string)),
			),
		}
	}

	return ghost.Result{
		Ok: false,
		Message: fmt.Sprintf(`%v != %v
want: %v
got:  %v
`, argWant, argGot, want, got),
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

// Error asserts that an error is non-nil.
func Error(err error) ghost.Result {
	args := ghostlib.ArgsFromAST(err)
	argErr := args[0]

	if err == nil {
		return ghost.Result{
			Ok:      false,
			Message: fmt.Sprintf("%s is nil", argErr),
		}
	}

	return ghost.Result{
		Ok:      true,
		Message: fmt.Sprintf("%s has error value: %s", argErr, err),
	}
}

// ErrorContaining asserts that an error string contains a particular substring.
func ErrorContaining(msg string, err error) ghost.Result {
	args := ghostlib.ArgsFromAST(msg, err)
	argMsg, argErr := args[0], args[1]

	switch {
	case err == nil && argMsg == fmt.Sprintf("%q", msg):
		return ghost.Result{
			Ok:      false,
			Message: fmt.Sprintf(`%v is nil; missing error message: %v`, argErr, msg),
		}
	case err == nil:
		return ghost.Result{
			Ok:      false,
			Message: fmt.Sprintf(`%v is nil; missing error message %v: %v`, argErr, argMsg, msg),
		}
	case strings.Contains(err.Error(), msg):
		return ghost.Result{
			Ok:      true,
			Message: fmt.Sprintf("%v contains error message %q: %v", argErr, msg, err),
		}
	default:
		return ghost.Result{
			Ok:      false,
			Message: fmt.Sprintf("%v does not contain error message %q: %v", argErr, msg, err),
		}
	}
}

// ErrorEqual asserts that an error string equals a particular message.
func ErrorEqual(msg string, err error) ghost.Result {
	args := ghostlib.ArgsFromAST(msg, err)
	argErr := args[1]

	if err == nil {
		return ghost.Result{
			Ok:      false,
			Message: fmt.Sprintf(`%v is nil; want message: %v`, argErr, msg),
		}
	}

	if err.Error() == msg {
		return ghost.Result{
			Ok:      true,
			Message: fmt.Sprintf("%v equals error message %q: %v", argErr, msg, err),
		}
	}

	return ghost.Result{
		Ok:      false,
		Message: fmt.Sprintf("%v does not equal error message %q: %v", argErr, msg, err),
	}
}

// False asserts that a value is false.
func False(b bool) ghost.Result {
	args := ghostlib.ArgsFromAST(b)
	argB := args[0]

	return ghost.Result{
		Ok:      !b,
		Message: fmt.Sprintf("%v is %t", argB, b),
	}
}

var jsonCompareOpts = jsondiff.DefaultConsoleOptions()

// JSONEqual asserts that two sets of JSON-encoded data are equivalent.
func JSONEqual[T ~string | ~[]byte](want, got T) ghost.Result {
	args := ghostlib.ArgsFromAST(want, got)
	argWant, argGot := args[0], args[1]

	diff, desc := jsondiff.Compare([]byte(want), []byte(got), &jsonCompareOpts)

	switch diff {
	case jsondiff.FullMatch:
		return ghost.Result{
			Ok:      true,
			Message: fmt.Sprintf("%v and %v are JSON equal", argWant, argGot),
		}
	case jsondiff.FirstArgIsInvalidJson:
		return ghost.Result{
			Ok: false,
			Message: fmt.Sprintf(`%v is not valid JSON
value: %s`, argWant, want),
		}
	case jsondiff.SecondArgIsInvalidJson:
		return ghost.Result{
			Ok: false,
			Message: fmt.Sprintf(`%v is not valid JSON
value: %s`, argGot, got),
		}
	case jsondiff.BothArgsAreInvalidJson:
		return ghost.Result{
			Ok: false,
			Message: fmt.Sprintf(`%v and %v are not valid JSON
want:
%s

got:
%s`, argWant, argGot, want, got),
		}
	}

	return ghost.Result{
		Ok: false,
		Message: fmt.Sprintf(`%v and %v are not JSON equal
diff: %s`, argWant, argGot, desc),
	}
}

// MapLen asserts that the length of a map is a particular size.
func MapLen[K comparable, V any](want int, got map[K]V) ghost.Result {
	args := ghostlib.ArgsFromAST(want, got)
	argGot := args[1]

	return ghost.Result{
		Ok: want == len(got),
		Message: fmt.Sprintf(`want %v length %d, got %d
map: %v
`, argGot, want, len(got), mapToString(got)),
	}
}

// mapToString pretty prints a map.
func mapToString[K comparable, V any](m map[K]V) string {
	var sb strings.Builder
	sb.WriteString("{\n")
	for k, v := range m {
		sb.WriteByte('\t')
		fmt.Fprintf(&sb, "%v: %v", k, v)
		sb.WriteString(",\n")
	}
	sb.WriteString("}")
	return sb.String()
}

// Nil asserts that the given value is nil.
func Nil(v any) ghost.Result {
	args := ghostlib.ArgsFromAST(v)
	argV := args[0]

	if isNil(v) {
		return ghost.Result{
			Ok:      true,
			Message: fmt.Sprintf("%v is nil", argV),
		}
	}

	return ghost.Result{
		Ok:      false,
		Message: fmt.Sprintf("%v is %v, not nil", argV, v),
	}
}

func isNil(v any) bool {
	if v == nil {
		return true
	}

	// Try reflection to catch typed nils
	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Pointer,
		reflect.Slice,
		reflect.UnsafePointer:
		return value.IsNil()
	}

	return false
}

// Panic asserts that the given function panics when invoked.
func Panic(f func()) (result ghost.Result) {
	args := ghostlib.ArgsFromAST(f)
	argF := args[0]

	defer func() {
		if r := recover(); r != nil {
			if strings.Contains(argF, "\n") {
				result = ghost.Result{
					Ok: true,
					Message: fmt.Sprintf(`function panicked with value: %v
%v
`, r, argF),
				}
			} else {
				result = ghost.Result{
					Ok:      true,
					Message: fmt.Sprintf(`function %v panicked with value: %v`, argF, r),
				}
			}
		}
	}()

	f()

	if strings.Contains(argF, "\n") {
		return ghost.Result{
			Ok: false,
			Message: fmt.Sprintf(`function did not panic
%v
`, argF),
		}
	}

	return ghost.Result{
		Ok:      false,
		Message: fmt.Sprintf("function %v did not panic", argF),
	}
}

// SliceContaining asserts that an element exists in a given slice.
func SliceContaining[T comparable](element T, slice []T) ghost.Result {
	args := ghostlib.ArgsFromAST(element, slice)
	argElement, argSlice := args[0], args[1]

	for _, x := range slice {
		if x == element {
			return ghost.Result{
				Ok: true,
				Message: fmt.Sprintf(`%v contains %v
element: %v
slice:   %v
`,
					argSlice,
					argElement,
					element,
					sliceElementToString(slice, element),
				),
			}
		}
	}

	return ghost.Result{
		Ok: false,
		Message: fmt.Sprintf(`%v does not contain %v
element: %v
slice:   %v
`,
			argSlice,
			argElement,
			element,
			sliceElementToString(slice, element),
		),
	}
}

// sliceElementToString pretty prints a slice, highlighting an element if it exists.
func sliceElementToString[T comparable](slice []T, element T) string {
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

// SliceLen asserts that the length of a slice is a particular size.
func SliceLen[T any](want int, got []T) ghost.Result {
	args := ghostlib.ArgsFromAST(want, got)
	argGot := args[1]

	return ghost.Result{
		Ok: want == len(got),
		Message: fmt.Sprintf(`want %v length %d, got %d
slice: %v
`, argGot, want, len(got), sliceToString(got)),
	}
}

// sliceToString pretty prints a slice.
func sliceToString[T any](slice []T) string {
	if len(slice) <= 3 {
		return fmt.Sprint(slice)
	}

	var sb strings.Builder
	sb.WriteString("[\n")
	for _, e := range slice {
		sb.WriteByte('\t')
		fmt.Fprint(&sb, e)
		sb.WriteByte('\n')
	}
	sb.WriteString("]")
	return sb.String()
}

// StringContaining asserts that a substring exists in a given string.
func StringContaining(substr, str string) ghost.Result {
	args := ghostlib.ArgsFromAST(substr, str)
	argSubstr, argStr := args[0], args[1]

	if strings.Contains(str, substr) {
		return ghost.Result{
			Ok: true,
			Message: fmt.Sprintf(`%v contains %v
substr: %s
str:    %s
`, argStr, argSubstr, quoteString(substr), quoteString(str)),
		}
	}

	return ghost.Result{
		Ok: false,
		Message: fmt.Sprintf(`%v does not contain %v
substr: %s
str:    %s
`, argStr, argSubstr, quoteString(substr), quoteString(str)),
	}
}

// True asserts that a value is true.
func True(b bool) ghost.Result {
	args := ghostlib.ArgsFromAST(b)
	argB := args[0]

	return ghost.Result{
		Ok:      b,
		Message: fmt.Sprintf("%v is %t", argB, b),
	}
}

// Zero asserts that the given value equals its zero value.
func Zero[T comparable](v T) ghost.Result {
	args := ghostlib.ArgsFromAST(v)
	argV := args[0]

	var zero T
	if v == zero {
		return ghost.Result{
			Ok:      true,
			Message: fmt.Sprintf("%v is the zero value", argV),
		}
	}

	if argV != fmt.Sprint(v) {
		return ghost.Result{
			Ok:      false,
			Message: fmt.Sprintf("%v is non-zero\nvalue: %v", argV, v),
		}
	}

	return ghost.Result{
		Ok:      false,
		Message: fmt.Sprintf("%v is non-zero", argV),
	}
}
