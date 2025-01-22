package errors

import (
	"fmt"
	"time"
)

type CustomError struct {
	StatusCode int
	Message    string
	Details    map[string]interface{}
	Timestamp  time.Time
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("Error code: %d, Message: %s, Details: %v, Timestamp: %s", e.StatusCode, e.Message, e.Details, e.Timestamp)
}

func New(code int, message string, details map[string]interface{}) *CustomError {
	return &CustomError{
		StatusCode: code,
		Message:    message,
		Details:    details,
		Timestamp:  time.Now(),
	}
}
