package be

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/rliebz/ghost"
	"github.com/rliebz/ghost/ghostlib"
)

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
func ErrorContaining(err error, msg string) ghost.Result {
	args := ghostlib.ArgsFromAST(err, msg)
	argErr, argMsg := args[0], args[1]

	switch {
	case err == nil && argMsg == fmt.Sprintf("%q", msg):
		return ghost.Result{
			Ok: false,
			Message: fmt.Sprintf(`error %v is nil, does not contain message
got:  <nil>
want: %v`,
				argErr,
				msg,
			),
		}
	case err == nil:
		return ghost.Result{
			Ok: false,
			Message: fmt.Sprintf(`error %v is nil, does not contain %v
got:  <nil>
want: %v`,
				argErr,
				argMsg,
				msg,
			),
		}
	case strings.Contains(err.Error(), msg):
		return ghost.Result{
			Ok: true,
			Message: fmt.Sprintf(`error %v contains message %v
got:  %v
want: %v`,
				argErr,
				argMsg,
				err,
				msg,
			),
		}
	default:
		return ghost.Result{
			Ok: false,
			Message: fmt.Sprintf(`error %v does not contain message %v
got:  %v
want: %v`,
				argErr,
				argMsg,
				err,
				msg,
			),
		}
	}
}

// ErrorEqual asserts that an error string equals a particular message.
func ErrorEqual(err error, msg string) ghost.Result {
	args := ghostlib.ArgsFromAST(err, msg)
	argErr, argMsg := args[0], args[1]

	if err == nil {
		return ghost.Result{
			Ok: false,
			Message: fmt.Sprintf(`error %v is nil
got:  <nil>
want: %v`,
				argErr,
				msg,
			),
		}
	}

	if err.Error() == msg {
		return ghost.Result{
			Ok: true,
			Message: fmt.Sprintf(`error %v has message %v
value: %v`,
				argErr,
				argMsg,
				err,
			),
		}
	}

	return ghost.Result{
		Ok: false,
		Message: fmt.Sprintf(`error %v does not have message %v
got:  %v
want: %v`,
			argErr,
			argMsg,
			err,
			msg,
		),
	}
}

// ErrorIs asserts that an error matches another using [errors.Is].
func ErrorIs(err error, target error) ghost.Result {
	args := ghostlib.ArgsFromAST(err, target)
	argErr, argTarget := args[0], args[1]

	if errors.Is(err, target) {
		return ghost.Result{
			Ok: true,
			Message: fmt.Sprintf(`error %v is target %v
error:  %v
target: %v`,
				argErr,
				argTarget,
				err,
				target,
			),
		}
	}

	return ghost.Result{
		Ok: false,
		Message: fmt.Sprintf(`error %v is not target %v
error:  %v
target: %v`,
			argErr,
			argTarget,
			err,
			target,
		),
	}
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()

// ErrorAs asserts that an error matches another using [errors.As].
func ErrorAs(err error, target any) ghost.Result {
	args := ghostlib.ArgsFromAST(err, target)
	argErr, argTarget := args[0], args[1]

	if err == nil {
		return ghost.Result{
			Ok:      false,
			Message: fmt.Sprintf("error %v was nil", argErr),
		}
	}

	// These next few checks are for invalid usage, where errors.As will panic if
	// a caller hits any of them. As an assertion library, it's probably more
	// polite to never panic.

	if target == nil {
		return ghost.Result{
			Ok:      false,
			Message: fmt.Sprintf("target %v cannot be nil", argTarget),
		}
	}

	val := reflect.ValueOf(target)
	typ := val.Type()
	if typ.Kind() != reflect.Ptr || val.IsNil() {
		return ghost.Result{
			Ok:      false,
			Message: fmt.Sprintf("target %v must be a non-nil pointer", argTarget),
		}
	}
	targetType := typ.Elem()

	if targetType.Kind() != reflect.Interface && !targetType.Implements(errorType) {
		return ghost.Result{
			Ok:      false,
			Message: fmt.Sprintf("*target %v must be interface or implement error", argTarget),
		}
	}

	if errors.As(err, target) {
		return ghost.Result{
			Ok: true,
			Message: fmt.Sprintf(`error %v set as target %v
error:  %v
target: %v`,
				argErr,
				argTarget,
				err,
				targetType,
			),
		}
	}

	return ghost.Result{
		Ok: false,
		Message: fmt.Sprintf(`error %v cannot be set as target %v
error:  %v
target: %v`,
			argErr,
			argTarget,
			err,
			targetType,
		),
	}
}
