package errors

import (
	"fmt"
)

var (
	ErrUsernameExists   = fmt.Errorf("username already exists")
	ErrUsernameNotFound = fmt.Errorf("username not found")
)

type CustomError struct {
	Code    int
	Message string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

func ItemNotFoundError(id int) *CustomError {
	return &CustomError{
		Message: fmt.Sprintf("Item not found. ID: %d", id),
	}
}
