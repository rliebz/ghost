package be

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rliebz/ghost"
	"github.com/rliebz/ghost/ghostlib"
)

// Error asserts that an error is non-nil.
func Error(err error) ghost.Result {
	args := ghostlib.ArgsFromAST(err)
	if err == nil {
		return ghost.Result{
			Ok: false,
			Message: func() string {
				return args.Get(0) + " is nil"
			},
		}
	}

	return ghost.Result{
		Ok: true,
		Message: func() string {
			return fmt.Sprintf("%s has error value: %s", args.Get(0), err)
		},
	}
}

// ErrorContaining asserts that an error string contains a particular substring.
func ErrorContaining(err error, msg string) ghost.Result {
	args := ghostlib.ArgsFromAST(err, msg)
	switch {
	case err == nil:
		return ghost.Result{
			Ok: false,
			Message: func() string {
				argErr, argMsg := args.Get(0), args.Get(1)
				if argMsg == fmt.Sprintf("%q", msg) {
					return fmt.Sprintf(`error %v is nil, does not contain message
got:  <nil>
want: %v`,
						argErr,
						msg,
					)
				}
				return fmt.Sprintf(`error %v is nil, does not contain %v
got:  <nil>
want: %v`,
					argErr,
					argMsg,
					msg,
				)
			},
		}
	case strings.Contains(err.Error(), msg):
		return ghost.Result{
			Ok: true,
			Message: func() string {
				argErr, argMsg := args.Get(0), args.Get(1)
				return fmt.Sprintf(`error %v contains message %v
got:  %v
want: %v`,
					argErr,
					argMsg,
					err,
					msg,
				)
			},
		}
	default:
		return ghost.Result{
			Ok: false,
			Message: func() string {
				argErr, argMsg := args.Get(0), args.Get(1)
				return fmt.Sprintf(`error %v does not contain message %v
got:  %v
want: %v`,
					argErr,
					argMsg,
					err,
					msg,
				)
			},
		}
	}
}

// ErrorEqual asserts that an error string equals a particular message.
func ErrorEqual(err error, msg string) ghost.Result {
	args := ghostlib.ArgsFromAST(err, msg)
	if err == nil {
		return ghost.Result{
			Ok: false,
			Message: func() string {
				return fmt.Sprintf(`error %v is nil
got:  <nil>
want: %v`,
					args.Get(0),
					msg,
				)
			},
		}
	}

	if err.Error() == msg {
		return ghost.Result{
			Ok: true,
			Message: func() string {
				return fmt.Sprintf(`error %v has message %v
value: %v`,
					args.Get(0),
					args.Get(1),
					err,
				)
			},
		}
	}

	return ghost.Result{
		Ok: false,
		Message: func() string {
			return fmt.Sprintf(`error %v does not have message %v
got:  %v
want: %v`,
				args.Get(0),
				args.Get(1),
				err,
				msg,
			)
		},
	}
}

// ErrorIs asserts that an error matches another using [errors.Is].
func ErrorIs(err error, target error) ghost.Result {
	args := ghostlib.ArgsFromAST(err, target)
	if errors.Is(err, target) {
		return ghost.Result{
			Ok: true,
			Message: func() string {
				return fmt.Sprintf(`error %v is target %v
error:  %v
target: %v`,
					args.Get(0),
					args.Get(1),
					err,
					target,
				)
			},
		}
	}

	return ghost.Result{
		Ok: false,
		Message: func() string {
			return fmt.Sprintf(`error %v is not target %v
error:  %v
target: %v`,
				args.Get(0),
				args.Get(1),
				err,
				target,
			)
		},
	}
}

// ErrorAs asserts that an error matches another using [errors.As].
func ErrorAs[T any](err error, target *T) ghost.Result {
	args := ghostlib.ArgsFromAST(err, target)
	if err == nil {
		return ghost.Result{
			Ok: false,
			Message: func() string {
				return fmt.Sprintf("error %v was nil", args.Get(0))
			},
		}
	}

	if target == nil {
		return ghost.Result{
			Ok: false,
			Message: func() string {
				return fmt.Sprintf("target %v cannot be nil", args.Get(1))
			},
		}
	}

	if errors.As(err, target) {
		return ghost.Result{
			Ok: true,
			Message: func() string {
				return fmt.Sprintf(`error %v set as target %v
error:  %v
target: %T`,
					args.Get(0),
					args.Get(1),
					err,
					*target,
				)
			},
		}
	}

	return ghost.Result{
		Ok: false,
		Message: func() string {
			return fmt.Sprintf(`error %v cannot be set as target %v
error:  %v
target: %T`,
				args.Get(0),
				args.Get(1),
				err,
				*target,
			)
		},
	}
}
