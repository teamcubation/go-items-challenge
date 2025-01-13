package application

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/teamcubation/go-items-challenge/internal/domain/item"
// 	"github.com/teamcubation/go-items-challenge/internal/ports/out/mocks"
// )

// func TestItemService_CreateItem(t *testing.T) {
// 	mockRepo := new(mocks.ItemRepository)
// 	service := NewItemService(mockRepo)

// 	newItem := &item.Item{
// 		Code:  "123",
// 		Stock: 10,
// 	}
// 	createdItem := &item.Item{
// 		ID:        1,
// 		Code:      "123",
// 		Stock:     10,
// 		Status:    "ACTIVE",
// 		CreatedAt: time.Now(),
// 		UpdatedAt: time.Now(),
// 	}

// 	mockRepo.On("ItemExistsByCode", mock.Anything, newItem.Code).Return(false)
// 	mockRepo.On("CreateItem", mock.Anything, newItem).Return(createdItem, nil)

// 	result, err := service.CreateItem(context.Background(), newItem)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, result)
// 	assert.Equal(t, createdItem, result)
// 	mockRepo.AssertCalled(t, "ItemExistsByCode", mock.Anything, newItem.Code)
// 	mockRepo.AssertCalled(t, "CreateItem", mock.Anything, newItem)
// }

// func TestItemService_GetItemById(t *testing.T) {
// 	mockRepo := new(mocks.ItemRepository)
// 	service := NewItemService(mockRepo)

// 	itemID := 1
// 	expectedItem := &item.Item{
// 		ID:    itemID,
// 		Code:  "123",
// 		Stock: 10,
// 	}

// 	mockRepo.On("GetItemById", mock.Anything, itemID).Return(expectedItem, nil)

// 	result, err := service.GetItemById(context.Background(), itemID)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, result)
// 	assert.Equal(t, expectedItem, result)
// 	mockRepo.AssertCalled(t, "GetItemById", mock.Anything, itemID)
// }

// func TestItemService_UpdateItem(t *testing.T) {
// 	mockRepo := new(mocks.ItemRepository)
// 	service := NewItemService(mockRepo)

// 	updatedItem := &item.Item{
// 		ID:    1,
// 		Code:  "123",
// 		Stock: 5,
// 	}

// 	existingItem := &item.Item{
// 		ID:        1,
// 		Code:      "123",
// 		Stock:     10,
// 		CreatedAt: time.Now(),
// 	}

// 	mockRepo.On("GetItemById", mock.Anything, updatedItem.ID).Return(existingItem, nil)
// 	mockRepo.On("UpdateItem", mock.Anything, updatedItem).Return(updatedItem, nil)

// 	result, err := service.UpdateItem(context.Background(), updatedItem)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, result)
// 	assert.Equal(t, updatedItem, result)
// 	mockRepo.AssertCalled(t, "GetItemById", mock.Anything, updatedItem.ID)
// 	mockRepo.AssertCalled(t, "UpdateItem", mock.Anything, updatedItem)
// }

// func TestItemService_DeleteItem(t *testing.T) {
// 	mockRepo := new(mocks.ItemRepository)
// 	service := NewItemService(mockRepo)

// 	itemID := 1
// 	deletedItem := &item.Item{
// 		ID:    itemID,
// 		Code:  "123",
// 		Stock: 10,
// 	}

// 	mockRepo.On("DeleteItem", mock.Anything, itemID).Return(deletedItem, nil)

// 	result, err := service.DeleteItem(context.Background(), itemID)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, result)
// 	assert.Equal(t, deletedItem, result)
// 	mockRepo.AssertCalled(t, "DeleteItem", mock.Anything, itemID)
// }

// func TestItemService_ListItems(t *testing.T) {
// 	mockRepo := new(mocks.ItemRepository)
// 	service := NewItemService(mockRepo)

// 	status := "ACTIVE"
// 	limit := 10
// 	page := 1
// 	items := []*item.Item{
// 		{ID: 1, Code: "123", Stock: 10, Status: "ACTIVE"},
// 		{ID: 2, Code: "456", Stock: 5, Status: "ACTIVE"},
// 	}

// 	mockRepo.On("ListItems", mock.Anything, status, limit, page).Return(items, nil)

// 	result, totalPages, err := service.ListItems(context.Background(), status, limit, page)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, result)
// 	assert.Equal(t, items, result)
// 	assert.Equal(t, 1, totalPages)
// 	mockRepo.AssertCalled(t, "ListItems", mock.Anything, status, limit, page)
// }
