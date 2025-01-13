package out

import "context"

type CategoryClient interface {
	IsAValidCategory(ctx context.Context, id int) (bool, error)
}
