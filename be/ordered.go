package be

import (
	"fmt"

	"github.com/rliebz/ghost"
	"github.com/rliebz/ghost/ghostlib"
	"github.com/rliebz/ghost/internal/constraints"
)

// Greater asserts that the first value provided is strictly greater than the second.
func Greater[T constraints.Ordered](a, b T) ghost.Result {
	args := ghostlib.ArgsFromAST(a, b)
	argA, argB := args[0], args[1]

	if a > b {
		return ghost.Result{
			Ok: true,
			Message: fmt.Sprintf(`%v is greater than %v`,
				inline(a, argA),
				inline(b, argB),
			),
		}
	}

	return ghost.Result{
		Ok: false,
		Message: fmt.Sprintf(`%v is not greater than %v`,
			inline(a, argA),
			inline(b, argB),
		),
	}
}

// Less asserts that the first value provided is strictly less than the second.
func Less[T constraints.Ordered](a, b T) ghost.Result {
	args := ghostlib.ArgsFromAST(a, b)
	argA, argB := args[0], args[1]

	if a < b {
		return ghost.Result{
			Ok: true,
			Message: fmt.Sprintf(`%v is less than %v`,
				inline(a, argA),
				inline(b, argB),
			),
		}
	}

	return ghost.Result{
		Ok: false,
		Message: fmt.Sprintf(`%v is not less than %v`,
			inline(a, argA),
			inline(b, argB),
		),
	}
}

func inline(val any, arg string) string {
	switch val := val.(type) {
	case string:
		if fmt.Sprintf("%q", val) == arg {
			return arg
		}
		return fmt.Sprintf("%v (%q)", arg, val)
	default:
		if fmt.Sprint(val) == arg {
			return arg
		}
		return fmt.Sprintf("%v (%v)", arg, val)
	}
}
