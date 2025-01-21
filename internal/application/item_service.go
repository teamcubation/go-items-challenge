package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/teamcubation/go-items-challenge/internal/domain/item"
	"github.com/teamcubation/go-items-challenge/internal/ports/out"
	"github.com/teamcubation/go-items-challenge/pkg/log"
)

type itemService struct {
	repo   out.ItemRepository
	client out.CategoryClient
}

func NewItemService(repo out.ItemRepository, client out.CategoryClient) *itemService {
	return &itemService{repo: repo, client: client}
}

func (s *itemService) CreateItem(ctx context.Context, item *item.Item) (*item.Item, error) {
	if item.Code == "" {
		return nil, errors.New("invalid request body")
	}

	// calling the client to validate the category
	isValid, err := s.client.IsAValidCategory(ctx, item.CategoryID)
	if err != nil {
		return nil, errors.New("client error")
	}
	if !isValid {
		return nil, errors.New("invalid Category")
	}

	if s.repo.ItemExistsByCode(ctx, item.Code) {
		return nil, errors.New("item with this code already exists")
	}
	item.ID = generateID()
	item.Status = determineStatus(item.Stock)
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	return s.repo.CreateItem(ctx, item)
}

func (s *itemService) GetItemByID(ctx context.Context, id int) (*item.Item, error) {
	logger := log.GetFromContext(ctx)
	logger.Info("Entering ItemService: GetItemById()")

	itm, err := s.repo.GetItemByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if itm == nil {
		return nil, fmt.Errorf("item with ID %d not found", id)
	}
	return itm, nil
}

func (s *itemService) UpdateItem(ctx context.Context, updatedItem *item.Item) (*item.Item, error) {
	existingItem, err := s.repo.GetItemByID(ctx, updatedItem.ID)
	if err != nil {
		return nil, err
	}
	if existingItem == nil {
		return nil, errors.New("item not found")
	}

	// Retain original values if new values are not provided
	if updatedItem.Title == "" {
		updatedItem.Title = existingItem.Title
	}
	if updatedItem.Description == "" {
		updatedItem.Description = existingItem.Description
	}
	if updatedItem.Price == 0 {
		updatedItem.Price = existingItem.Price
	}

	updatedItem.CreatedAt = existingItem.CreatedAt
	updatedItem.UpdatedAt = time.Now()

	result, err := s.repo.UpdateItem(ctx, updatedItem)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *itemService) DeleteItem(ctx context.Context, id int) (*item.Item, error) {
	return s.repo.DeleteItem(ctx, id)
}

func (s *itemService) ListItems(ctx context.Context, status string, limit int, page int) ([]*item.Item, int, error) {
	if status == "" {
		status = "ACTIVE"
	}

	items, err := s.repo.ListItems(ctx, status, limit, page)
	if err != nil {
		return nil, 0, err
	}

	totalPages := (len(items.Data) + limit - 1) / limit
	var result []*item.Item
	for _, itm := range items.Data {
		result = append(result, &itm)
	}
	return result, totalPages, nil
}

func (s *itemService) ItemExistsByCode(_ context.Context, _ string) bool {
	return true
}

func generateID() int {
	return 0
}

func determineStatus(stock int) string {
	if stock > 0 {
		return "ACTIVE"
	}
	return "INACTIVE"
}
