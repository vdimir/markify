package app

import (
	"fmt"
)

// UserError error caused by wrong user input.
// Contains causer error and message intended to show to user
type UserError struct {
	Inner   error
	UserMsg string
}

// WrapfUserError return UserError with formatted user message
func WrapfUserError(err error, format string, a ...interface{}) UserError {
	return UserError{
		Inner:   err,
		UserMsg: fmt.Sprintf(format, a...),
	}
}

func (e UserError) String() string {
	return e.UserMsg
}

func (e UserError) Error() string {
	return fmt.Sprintf("UserError: %s, Msg: %s", e.Inner, e.UserMsg)
}
