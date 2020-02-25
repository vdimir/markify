package apperr

import (
	"fmt"
)

// UserError error caused by wrong user input
type UserError struct {
	Inner error
}

func (e UserError) Error() string {
	return fmt.Sprintf("UserError: %s", e.Inner)
}

// DBError error occurred in database operation
type DBError struct {
	Inner error
}

func (e DBError) Error() string {
	return fmt.Sprintf("DBError: %s", e.Inner)
}
