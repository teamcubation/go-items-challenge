package user

import (
	"github.com/teamcubation/go-items-challenge/internal/adapters/http"
	"github.com/teamcubation/go-items-challenge/internal/adapters/repository"
	"github.com/teamcubation/go-items-challenge/internal/application"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(repository.NewElasticsearchUserAdapter),
	fx.Provide(application.NewAuthService),
	fx.Provide(http.NewAuthHandler),
)
