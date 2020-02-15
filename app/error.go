package app

import (
	"fmt"
)

// UserError error caused by wrong user input
type UserError struct {
	inner error
}

func (e UserError) Error() string {
	return fmt.Sprintf("UserError: %s", e.inner)
}

// DBError error occured in database operation
type DBError struct {
	inner error
}

func (e DBError) Error() string {
	return fmt.Sprintf("DBError: %s", e.inner)
}
