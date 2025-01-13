package repository

import (
	"context"
	"os"
	"testing"

	"gorm.io/driver/postgres"

	"github.com/stretchr/testify/assert"
	"github.com/teamcubation/go-items-challenge/internal/domain/item"

	"gorm.io/gorm"
)

func setupTestDB() (*gorm.DB, error) {
	dsn := os.Getenv("DSN")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&ItemModel{})
	return db, nil
}

func TestItemRepository(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}
	repo := NewItemRepository(db)

	t.Run("CreateItem", func(t *testing.T) {
		ctx := context.Background()
		newItem := &item.Item{
			Code:  "123",
			Stock: 10,
		}
		createdItem, err := repo.CreateItem(ctx, newItem)
		assert.NoError(t, err)
		assert.NotNil(t, createdItem)
		assert.Equal(t, newItem.Code, createdItem.Code)
	})

	t.Run("GetItemById", func(t *testing.T) {
		ctx := context.Background()
		newItem := &item.Item{
			Code:  "456",
			Stock: 5,
		}
		createdItem, _ := repo.CreateItem(ctx, newItem)
		fetchedItem, err := repo.GetItemById(ctx, createdItem.ID)
		assert.NoError(t, err)
		assert.NotNil(t, fetchedItem)
		assert.Equal(t, createdItem.Code, fetchedItem.Code)
	})

	t.Run("UpdateItem", func(t *testing.T) {
		ctx := context.Background()
		newItem := &item.Item{
			Code:  "789",
			Stock: 15,
		}
		createdItem, _ := repo.CreateItem(ctx, newItem)
		createdItem.Stock = 20
		updatedItem, err := repo.UpdateItem(ctx, createdItem)
		assert.NoError(t, err)
		assert.NotNil(t, updatedItem)
		assert.Equal(t, 20, updatedItem.Stock)
	})

	t.Run("DeleteItem", func(t *testing.T) {
		ctx := context.Background()
		newItem := &item.Item{
			Code:  "101",
			Stock: 25,
		}
		createdItem, _ := repo.CreateItem(ctx, newItem)
		deletedItem, err := repo.DeleteItem(ctx, createdItem.ID)
		assert.NoError(t, err)
		assert.NotNil(t, deletedItem)
		assert.Equal(t, createdItem.Code, deletedItem.Code)
	})

	t.Run("ListItems", func(t *testing.T) {
		ctx := context.Background()
		status := "ACTIVE"
		limit := 10
		page := 1
		items, err := repo.ListItems(ctx, status, limit, page)
		assert.NoError(t, err)
		assert.NotNil(t, items)
	})
}
