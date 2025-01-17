package in

import (
	"context"

	"github.com/teamcubation/go-items-challenge/internal/domain/item"
)

type ItemService interface {
	CreateItem(ctx context.Context, itm *item.Item) (*item.Item, error)
	UpdateItem(ctx context.Context, itm *item.Item) (*item.Item, error)
	DeleteItem(ctx context.Context, id int) (*item.Item, error)
	GetItemByID(ctx context.Context, id int) (*item.Item, error)
	ListItems(ctx context.Context, status string, limit int, page int) ([]*item.Item, int, error)
	ItemExistsByCode(ctx context.Context, code string) bool
}
