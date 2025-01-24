package errors

import (
	"fmt"
)

var (
	ErrUsernameExists   = fmt.Errorf("username already exists")
	ErrUsernameNotFound = fmt.Errorf("username not found")
	ErrFetchingUser     = fmt.Errorf("error fetching user")
	ErrCreatingUser     = fmt.Errorf("error creating user")
	ErrHashingPassword  = fmt.Errorf("error hashing password")
	ErrTokenGeneration  = fmt.Errorf("error generating token")
	ErrItemNotFound     = fmt.Errorf("item not found")
	ErrUpdatingItem     = fmt.Errorf("error updating item")
	ErrDeletingItem     = fmt.Errorf("error deleting item")
	ErrCreatingItem     = fmt.Errorf("error creating item")
	ErrFetchingItem     = fmt.Errorf("error fetching item")
	ErrFetchingItems    = fmt.Errorf("error fetching items")
	ErrRequestBody      = fmt.Errorf("invalid request body")
	ErrClientError      = fmt.Errorf("error with client")
	ErrInvalidCategory  = fmt.Errorf("invalid category")
	ErrCodeExists       = fmt.Errorf("item with this code already exists")
	ErrInternalServer   = fmt.Errorf("internal server error")
	ErrEncodingResponse = fmt.Errorf("error encoding response")
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
