package item

import (
	"github.com/teamcubation/go-items-challenge/internal/adapters/client"
	"github.com/teamcubation/go-items-challenge/internal/adapters/http"
	"github.com/teamcubation/go-items-challenge/internal/adapters/repository"
	"github.com/teamcubation/go-items-challenge/internal/application"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(http.NewItemHandler),
	fx.Provide(repository.NewElasticsearchAdapter),
	fx.Provide(func() string {
		return "http://mockapi:8000"
	}),
	fx.Provide(client.NewCategoryClient),
	fx.Provide(application.NewItemService),
)
