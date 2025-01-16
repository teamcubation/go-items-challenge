package out

import (
	"context"

	"github.com/teamcubation/go-items-challenge/internal/domain/item"
	"github.com/teamcubation/go-items-challenge/internal/domain/user"
)

type ItemRepository interface {
	CreateItem(ctx context.Context, itm *item.Item) (*item.Item, error)
	GetItemByID(ctx context.Context, id int) (*item.Item, error)
	UpdateItem(ctx context.Context, itm *item.Item) (*item.Item, error)
	DeleteItem(ctx context.Context, id int) (*item.Item, error)
	ItemExistsByCode(ctx context.Context, code string) bool
	ListItems(ctx context.Context, status string, limit int, page int) (*item.Response, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, u *user.User) error
	GetUserByUsername(ctx context.Context, username string) (*user.User, error)
}
