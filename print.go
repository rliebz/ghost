package ghost

import (
	"fmt"
	"strings"
)

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
