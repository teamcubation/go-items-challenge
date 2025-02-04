package modules

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/teamcubation/go-items-challenge/cmd/server/modules/item"
	"github.com/teamcubation/go-items-challenge/cmd/server/modules/user"
	httphdl "github.com/teamcubation/go-items-challenge/internal/adapters/http"
	"github.com/teamcubation/go-items-challenge/internal/adapters/http/middleware"
	"go.uber.org/fx"
)

func NewApp(decorators ...any) *fx.App {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	options := []fx.Option{
		// fx.Provide() registers constructors for dependencies to be injected.
		fx.Provide(
			httphdl.NewExportHandler,
			NewRouter,
		),
		item.Module,
		user.Module,
	}

	// fx.New() creates a new application with registered dependencies and lifecycle hooks.
	return fx.New(
		// fx.Options() groups multiple fx.Option configurations into one.
		fx.Options(options...),
		// fx.Invoke() calls functions once dependencies are available.
		fx.Invoke(StartServer),
	)
}

func NewRouter(
	authHandler *httphdl.AuthHandler,
	itemHandler *httphdl.ItemHandler,
	exportHandler *httphdl.ExportHandler,
) *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)

	api.HandleFunc("/items/export", exportHandler.Export).Methods("GET")
	api.HandleFunc("/items", itemHandler.CreateItem).Methods("POST")
	api.HandleFunc("/items/{id}", itemHandler.UpdateItem).Methods("PUT")
	api.HandleFunc("/items/{id}", itemHandler.DeleteItem).Methods("DELETE")
	api.HandleFunc("/items/{id}", itemHandler.GetItemByID).Methods("GET")
	api.HandleFunc("/items", itemHandler.ListItems).Methods("GET")

	return r
}

func StartServer(lc fx.Lifecycle, router *mux.Router) *http.Server {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// fx.Lifecycle allows us to hook into the application startup and shutdown.
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("Error running server: %v", err)
				}
			}()
			log.Println("Server running on port 8080")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down server...")
			shutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			return srv.Shutdown(shutCtx)
		},
	})

	return srv
}
