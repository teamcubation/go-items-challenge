package repository

import (
	"time"

	"gorm.io/gorm"
)

type ItemModel struct {
	ID          uint      `gorm:"primaryKey"`
	Code        string    `json:"code"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Status      string    `json:"status" gorm:"-"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (ItemModel) TableName() string {
	return "items"
}

func (item *ItemModel) BeforeSave(_ *gorm.DB) (err error) {
	if item.Stock > 0 {
		item.Status = "ACTIVE"
	} else {
		item.Status = "INACTIVE"
	}
	return
}
