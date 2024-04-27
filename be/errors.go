package be

import (
	"fmt"
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
