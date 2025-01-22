package presenter

import (
	"time"
)

type ApiError struct {
	StatusCode int                    `json:"status_code"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details"`
	Timestamp  time.Time              `json:"timestamp"`
}

func NewApiError(code int, message string, details map[string]interface{}) ApiError {
	return ApiError{
		StatusCode: code,
		Message:    message,
		Details:    details,
		Timestamp:  time.Now(),
	}
}
