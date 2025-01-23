package server

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/teamcubation/go-items-challenge/internal/adapters/client"
	httphdl "github.com/teamcubation/go-items-challenge/internal/adapters/http"
	"github.com/teamcubation/go-items-challenge/internal/adapters/http/middleware"
	"github.com/teamcubation/go-items-challenge/internal/adapters/repository"
	"github.com/teamcubation/go-items-challenge/internal/application"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, externalAPIURL string) *mux.Router {
	userRepo := repository.NewUserRepository(db)
	userSrv := application.NewAuthService(userRepo)
	authHandler := httphdl.NewAuthHandler(userSrv)

	itemRepo := repository.NewItemRepository(db)
	categoryClient := client.NewCategoryClient(externalAPIURL)
	itemSrv := application.NewItemService(itemRepo, categoryClient)
	itemHandler := httphdl.NewItemHandler(itemSrv)
	exportHandler := httphdl.NewExportHandler(itemSrv)

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
