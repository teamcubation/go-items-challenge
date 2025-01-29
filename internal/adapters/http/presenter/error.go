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
