package be

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/rliebz/ghost"
	"github.com/rliebz/ghost/ghostlib"
	"github.com/rliebz/ghost/internal/constraints"
	"github.com/rliebz/ghost/internal/jsondiff"
)

// AssignedAs assigns a value to a target of an arbitrary type.
//
// The target must be a non-nil pointer.
func AssignedAs[T any](value any, target *T) ghost.Result {
	args := ghostlib.ArgsFromAST(value, target)
	argValue, argTarget := args[0], args[1]

	if target == nil {
		return ghost.Result{
			Ok:      false,
			Message: fmt.Sprintf("target %s cannot be nil", argTarget),
		}
	}

	typedValue, ok := value.(T)
	if !ok {
		return ghost.Result{
			Ok: false,
			Message: fmt.Sprintf(
				`%s (%T) could not be assigned to %s (%T)
value: %v`,
				argValue,
				value,
				argTarget,
				target,
				value,
			),
		}
	}

	*target = typedValue

	return ghost.Result{
		Ok: true,
		Message: fmt.Sprintf(
			`%s (%T) was assigned to %s (%T)
value: %v`,
			argValue,
			value,
			argTarget,
			target,
			value,
		),
	}
}

// Close asserts that a value is within a delta of another.
func Close[T constraints.Integer | constraints.Float](got, want, delta T) ghost.Result {
	args := ghostlib.ArgsFromAST(got, want, delta)
	argGot, argWant := args[0], args[1]

	gotDelta := want - got
	if gotDelta < 0 {
		gotDelta = 0 - gotDelta
	}

	if _, err := strconv.ParseFloat(argGot, 64); err != nil {
		argGot = fmt.Sprintf("%s (%v)", argGot, got)
	}

	if _, err := strconv.ParseFloat(argWant, 64); err != nil {
		argWant = fmt.Sprintf("%s (%v)", argWant, want)
	}

	if gotDelta <= delta {
		return ghost.Result{
			Ok: true,
			Message: fmt.Sprintf(
				`delta %v between %s and %s is within %v
got:   %v
want:  %v
delta: %v`,
				gotDelta, argGot, argWant, delta,
				got,
				want,
				gotDelta,
			),
		}
	}

	return ghost.Result{
		Ok: false,
		Message: fmt.Sprintf(
			`delta %v between %s and %s is not within %v
got:   %v
want:  %v
delta: %v`,
			gotDelta, argGot, argWant, delta,
			got,
			want,
			gotDelta,
		),
	}
}

// DeepEqual asserts that two elements are deeply equal.
func DeepEqual[T any](got, want T) ghost.Result {
	args := ghostlib.ArgsFromAST(got, want)
	argGot, argWant := args[0], args[1]

	if diff := colorDiff(want, got); diff != "" {
		return ghost.Result{
			Ok: false,
			Message: fmt.Sprintf(`%v != %v
%v`, argGot, argWant, diff),
		}
	}

	return ghost.Result{
		Ok: true,
		Message: fmt.Sprintf(`%v == %v
value: %v
`, argGot, argWant, want),
	}
}

// Equal asserts that two elements are equal.
func Equal[T comparable](got T, want T) ghost.Result {
	args := ghostlib.ArgsFromAST(got, want)
	argGot, argWant := args[0], args[1]

	if got == want {
		switch fmt.Sprint(want) {
		case argGot, argWant:
			return ghost.Result{
				Ok:      true,
				Message: fmt.Sprintf(`%v == %v`, argGot, argWant),
			}
		default:
			return ghost.Result{
				Ok: true,
				Message: fmt.Sprintf(`%v == %v
value: %v
`, argGot, argWant, want),
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
%v`, argGot, argWant, colorDiff(want, got)),
		}
	case reflect.String:
		if strings.ContainsAny(v.String(), "\n\r") ||
			strings.ContainsAny(reflect.ValueOf(got).String(), "\n\r") {

			return ghost.Result{
				Ok: false,
				Message: fmt.Sprintf(`%v != %v
%v`, argGot, argWant, colorDiff(want, got)),
			}
		}

		return ghost.Result{
			Ok: false,
			Message: fmt.Sprintf(`%v != %v
got:  %v
want: %v
`,
				argGot,
				argWant,
				quoteString(reflect.ValueOf(got).String()),
				quoteString(reflect.ValueOf(want).String()),
			),
		}
	}

	return ghost.Result{
		Ok: false,
		Message: fmt.Sprintf(`%v != %v
got:  %v
want: %v
`, argGot, argWant, got, want),
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

// False asserts that a value is false.
func False(b bool) ghost.Result {
	args := ghostlib.ArgsFromAST(b)
	argB := args[0]

	return ghost.Result{
		Ok:      !b,
		Message: fmt.Sprintf("%v is %t", argB, b),
	}
}

// JSONEqual asserts that two sets of JSON-encoded data are equivalent.
func JSONEqual[T ~string | ~[]byte](got, want T) ghost.Result {
	args := ghostlib.ArgsFromAST(got, want)
	argGot, argWant := args[0], args[1]

	diff, kind := colorJSONDiff(got, want)

	switch kind {
	case jsondiff.Match:
		return ghost.Result{
			Ok:      true,
			Message: fmt.Sprintf("%v and %v are JSON equal", argGot, argWant),
		}
	case jsondiff.GotInvalid:
		return ghost.Result{
			Ok: false,
			Message: fmt.Sprintf(`%v is not valid JSON
value: %s`, argGot, got),
		}
	case jsondiff.WantInvalid:
		return ghost.Result{
			Ok: false,
			Message: fmt.Sprintf(`%v is not valid JSON
value: %s`, argWant, want),
		}
	case jsondiff.BothInvalid:
		return ghost.Result{
			Ok: false,
			Message: fmt.Sprintf(`%v and %v are not valid JSON
got:
%s

want:
%s`, argGot, argWant, got, want),
		}
	}

	return ghost.Result{
		Ok: false,
		Message: fmt.Sprintf(`%v and %v are not JSON equal
%s`, argGot, argWant, diff),
	}
}

// MapLen asserts that the length of a map is a particular size.
func MapLen[K comparable, V any](got map[K]V, want int) ghost.Result {
	args := ghostlib.ArgsFromAST(got, want)
	argGot := args[0]

	if len(got) == want {
		return ghost.Result{
			Ok: true,
			Message: fmt.Sprintf(`%v is length %d
map: %v
`, argGot, len(got), mapToString(got)),
		}
	}

	return ghost.Result{
		Ok: false,
		Message: fmt.Sprintf(`%v is length %d, not %d
map: %v
`, argGot, len(got), want, mapToString(got)),
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

// SliceContaining asserts that an element exists in a given slice.
func SliceContaining[T comparable](slice []T, element T) ghost.Result {
	args := ghostlib.ArgsFromAST(slice, element)
	argSlice, argElement := args[0], args[1]

	for _, x := range slice {
		if x == element {
			return ghost.Result{
				Ok: true,
				Message: fmt.Sprintf(`%v contains %v
slice:   %v
element: %v
`,
					argSlice,
					argElement,
					sliceElementToString(slice, element),
					element,
				),
			}
		}
	}

	return ghost.Result{
		Ok: false,
		Message: fmt.Sprintf(`%v does not contain %v
slice:   %v
element: %v
`,
			argSlice,
			argElement,
			sliceElementToString(slice, element),
			element,
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
func SliceLen[T any](got []T, want int) ghost.Result {
	args := ghostlib.ArgsFromAST(got, want)
	argGot := args[0]

	if len(got) == want {
		return ghost.Result{
			Ok: true,
			Message: fmt.Sprintf(`%v is length %d
slice: %v
`, argGot, len(got), sliceToString(got)),
		}
	}

	return ghost.Result{
		Ok: false,
		Message: fmt.Sprintf(`%v is length %d, not %d
slice: %v
`, argGot, len(got), want, sliceToString(got)),
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
func StringContaining(str, substr string) ghost.Result {
	args := ghostlib.ArgsFromAST(str, substr)
	argStr, argSubstr := args[0], args[1]

	if strings.Contains(str, substr) {
		return ghost.Result{
			Ok: true,
			Message: fmt.Sprintf(`%v contains %v
str:    %s
substr: %s
`, argStr, argSubstr, quoteString(str), quoteString(substr)),
		}
	}

	return ghost.Result{
		Ok: false,
		Message: fmt.Sprintf(`%v does not contain %v
str:    %s
substr: %s
`, argStr, argSubstr, quoteString(str), quoteString(substr)),
	}
}

// StringMatching asserts that a given string matches a regular expression.
func StringMatching(str string, expr string) ghost.Result {
	args := ghostlib.ArgsFromAST(str, expr)
	argStr, argExpr := args[0], args[1]

	re, err := regexp.Compile(expr)
	if err != nil {
		return ghost.Result{
			Ok: false,
			Message: fmt.Sprintf(`%v is not a valid regular expression
%v
`,
				argExpr,
				err,
			),
		}
	}

	if re.MatchString(str) {
		return ghost.Result{
			Ok: true,
			Message: fmt.Sprintf(`%v matches regular expression %v
str:  %s
expr: %s
`,
				argStr, argExpr,
				quoteString(str),
				re.String(),
			),
		}
	}

	return ghost.Result{
		Ok: false,
		Message: fmt.Sprintf(`%v does not match regular expression %v
str:  %s
expr: %s
`,
			argStr, argExpr,
			quoteString(str),
			re.String(),
		),
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
