package presenter

import "time"

type CustomError struct {
	StatusCode int                    `json:"status_code"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details"`
	Timestamp  time.Time              `json:"timestamp"`
}

func NewErrorResponse(code int, message string, details map[string]interface{}) CustomError {
	return CustomError{
		StatusCode: code,
		Message:    message,
		Details:    details,
		Timestamp:  time.Now(),
	}
}
