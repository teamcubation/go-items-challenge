package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/teamcubation/go-items-challenge/internal/adapters/http/middleware"
	"github.com/teamcubation/go-items-challenge/internal/domain/user"
	"github.com/teamcubation/go-items-challenge/pkg/log"

	"github.com/teamcubation/go-items-challenge/internal/domain/item"
	"gorm.io/gorm"

	"strings"
)

type itemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) *itemRepository {
	return &itemRepository{db: db}
}

func (r *itemRepository) CreateItem(ctx context.Context, itm *item.Item) (*item.Item, error) {
	userID, ok := ctx.Value(middleware.UserContextKey).(int)
	if !ok || userID == 0 {
		return nil, fmt.Errorf("invalid user ID in context")
	}
	itm.CreatedBy = userID
	itm.UpdatedBy = userID

	// checking if the item code already exists
	if r.ItemExistsByCode(ctx, itm.Code) {
		return nil, fmt.Errorf("item with code %s already exists", itm.Code)
	}

	// setting the status of the item
	if itm.Stock > 0 {
		itm.Status = "ACTIVE"
	} else {
		itm.Status = "INACTIVE"
	}

	// checking if the user exists
	var use user.User
	if err := r.db.WithContext(ctx).First(&use, userID).Error; err != nil {
		return nil, fmt.Errorf("user with ID %d not found", userID)
	}

	if err := r.db.WithContext(ctx).Create(itm).Error; err != nil {
		return nil, err
	}

	return itm, nil
}

func (r *itemRepository) GetItemById(ctx context.Context, id int) (*item.Item, error) {
	logger := log.GetFromContext(ctx)
	logger.Info("Entering itemRepository: GetItemById()")
	logger.Printf("Fetching item with ID: %d", id)

	var itm item.Item
	if err := r.db.WithContext(ctx).First(&itm, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("item with ID %d not found: %v", id, err)
		}
		return nil, err
	}

	if itm.Stock == 0 {
		itm.Status = "INACTIVE"
	} else {
		itm.Status = "ACTIVE"
	}

	if err := r.db.WithContext(ctx).Save(&itm).Error; err != nil {
		return nil, err
	}

	return &itm, nil
}

func (r *itemRepository) UpdateItem(ctx context.Context, itm *item.Item) (*item.Item, error) {

	userID, ok := ctx.Value(middleware.UserContextKey).(int)
	if !ok || userID == 0 {
		return nil, fmt.Errorf("user ID not found in context")
	}

	var existingItem item.Item
	if err := r.db.WithContext(ctx).First(&existingItem, itm.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("item with ID %d not found", itm.ID)
		}
		return nil, err
	}

	// Verify if the fields created_by, created_at, or updated_by were changed
	if itm.CreatedBy != 0 && existingItem.CreatedBy != itm.CreatedBy {
		return nil, fmt.Errorf("cannot change the created_by field")
	}
	if !itm.CreatedAt.IsZero() && !existingItem.CreatedAt.Equal(itm.CreatedAt) {
		return nil, fmt.Errorf("cannot change the created_at field")
	}
	//if itm.UpdatedBy != 0 && existingItem.UpdatedBy != itm.UpdatedBy {
	//	return nil, fmt.Errorf("cannot change the updated_by field")
	//}

	// Ensure the code field is not empty
	if itm.Code == "" {
		return nil, fmt.Errorf("code field cannot be empty")
	}

	// Only allow update if the item code is the same
	if existingItem.Code != itm.Code {
		return nil, fmt.Errorf("cannot update item with different code")
	}

	existingItem.Title = itm.Title
	existingItem.Description = itm.Description
	existingItem.Price = itm.Price
	existingItem.Stock = itm.Stock
	existingItem.UpdatedAt = time.Now()
	existingItem.UpdatedBy = userID

	if existingItem.Stock == 0 {
		existingItem.Status = "INACTIVE"
	} else {
		existingItem.Status = "ACTIVE"
	}

	if err := r.db.WithContext(ctx).Save(&existingItem).Error; err != nil {
		return nil, err
	}

	itm.Status = existingItem.Status
	itm.CreatedBy = existingItem.CreatedBy
	itm.CreatedAt = existingItem.CreatedAt
	itm.UpdatedBy = existingItem.UpdatedBy
	itm.UpdatedAt = existingItem.UpdatedAt

	return itm, nil
}

func (r *itemRepository) DeleteItem(ctx context.Context, id int) (*item.Item, error) {
	var itm item.Item
	if err := r.db.WithContext(ctx).First(&itm, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("item with ID %d not found", id)
		}
		return nil, err
	}
	if err := r.db.WithContext(ctx).Delete(&itm).Error; err != nil {
		return nil, err
	}
	return &itm, nil
}

func (r *itemRepository) ListItems(ctx context.Context, status string, limit int, page int) (*item.ItemResponse, error) {
	status = strings.ToUpper(status)
	if status != "ACTIVE" && status != "INACTIVE" {
		return nil, fmt.Errorf("Invalid status: %s", status)
	}

	var items []item.Item
	offset := (page - 1) * limit
	result := r.db.WithContext(ctx).Where("UPPER(status) = ?", status).Limit(limit).Offset(offset).Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}

	var totalItems int64
	r.db.WithContext(ctx).Model(&item.Item{}).Where("UPPER(status) = ?", status).Count(&totalItems)
	totalPages := int((totalItems + int64(limit) - 1) / int64(limit))

	response := &item.ItemResponse{
		TotalPages: totalPages,
		Data:       items,
	}

	return response, nil

}

func (r *itemRepository) ItemExistsByCode(ctx context.Context, code string) bool {
	var count int64
	r.db.WithContext(ctx).Model(&item.Item{}).Where("code = ?", code).Count(&count)
	return count > 0
}
