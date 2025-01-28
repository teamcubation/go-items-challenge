// errors/errors.go

package presenter

import (
	"fmt"
	"time"
)

type CustomError struct {
	Code      string                 `json:"code"`      // Código único do erro
	Message   string                 `json:"message"`   // Mensagem descritiva
	Details   map[string]interface{} `json:"details"`   // Metadados adicionais
	Timestamp time.Time              `json:"timestamp"` // Timestamp do erro
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Função auxiliar para criar novos erros
func New(code string, message string, details map[string]interface{}) *CustomError {
	return &CustomError{
		Code:      code,
		Message:   message,
		Details:   details,
		Timestamp: time.Now(),
	}
}

// Define erros específicos com códigos únicos
var (
	ErrFetchingUser       = New("ERR_FETCHING_USER", "Error fetching user", nil)
	ErrHashingPassword    = New("ERR_HASHING_PASSWORD", "Error hashing password", nil)
	ErrUsernameExists     = New("ERR_USERNAME_EXISTS", "Username already exists", nil)
	ErrCreatingUser       = New("ERR_CREATING_USER", "Error creating user", nil)
	ErrUsernameNotFound   = New("ERR_USERNAME_NOT_FOUND", "Username not found", nil)
	ErrTokenGeneration    = New("ERR_TOKEN_GENERATION", "Error generating token", nil)
	ErrInvalidCredentials = New("ERR_INVALID_CREDENTIALS", "Invalid credentials", nil)
	ErrInvalidRequestBody = New("ERR_INVALID_REQUEST_BODY", "Invalid request body", nil)
	ErrInvalidCategory    = New("ERR_INVALID_CATEGORY", "Invalid category", nil)
	ErrCodeExists         = New("ERR_CODE_EXISTS", "Code already exists", nil)
	ErrFetchingItem       = New("ERR_FETCHING_ITEM", "Error fetching item(s)", nil)
	ErrItemNotFound       = New("ERR_ITEM_NOT_FOUND", "Item not found", nil)
	ErrUpdatingItem       = New("ERR_UPDATING_ITEM", "Error updating item", nil)
	ErrMissingFields      = New("ERR_MISSING_FIELDS", "Missing required fields", nil)
	ErrValidationFields   = New("ERR_VALIDATION_FIELDS", "Invalid fields", nil)
)
