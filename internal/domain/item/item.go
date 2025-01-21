package item

import (
	"time"
)

type Item struct {
	ID          int       `json:"id"`
	Code        string    `json:"code" validate:"required,alphanum"`
	Title       string    `json:"title,omitempty" validate:"omitempty,min=4"`
	Description string    `json:"description,omitempty" validate:"omitempty,max=255"`
	CategoryID  int       `json:"category_id"`
	Price       float64   `json:"price,omitempty" validate:"omitempty,gt=0"`
	Stock       int       `json:"stock,omitempty" validate:"omitempty,min=0"`
	Status      string    `json:"status,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   int       `json:"created_by"`
	UpdatedBy   int       `json:"updated_by"`
}
