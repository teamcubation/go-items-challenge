package main

import (
	"context"
	_ "encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/docker/docker/client"

	"github.com/teamcubation/go-items-challenge/internal/application"
	"github.com/teamcubation/go-items-challenge/internal/domain/user"
	_ "github.com/teamcubation/go-items-challenge/internal/ports/in"
	_ "github.com/teamcubation/go-items-challenge/internal/ports/out"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/teamcubation/go-items-challenge/docs"
	"github.com/teamcubation/go-items-challenge/internal/adapters/client"
	httphdl "github.com/teamcubation/go-items-challenge/internal/adapters/http"
	"github.com/teamcubation/go-items-challenge/internal/adapters/http/middleware"
	"github.com/teamcubation/go-items-challenge/internal/adapters/repository"
	"github.com/teamcubation/go-items-challenge/internal/domain/item"

	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func runMigrations(db *gorm.DB) {
	err := db.AutoMigrate(&user.User{}, &item.Item{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}

// @title Items API
// @version 1.0
// @description This is a simple items API
// @host localhost:8080
// @BasePath /api
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	runMigrations(db)

	userRepo := repository.NewUserRepository(db)
	userSrv := application.NewAuthService(userRepo)
	authHandler := httphdl.NewAuthHandler(userSrv)

	itemRepo := repository.NewItemRepository(db)
	categoryClient := client.NewCategoryClient("http://localhost:8000")
	itemSrv := application.NewItemService(itemRepo, categoryClient)
	itemHandler := httphdl.NewItemHandler(itemSrv)

	r := mux.NewRouter()

	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)

	api.HandleFunc("/items", itemHandler.CreateItem).Methods("POST")
	api.HandleFunc("/items/{id}", itemHandler.UpdateItem).Methods("PUT")
	api.HandleFunc("/items/{id}", itemHandler.DeleteItem).Methods("DELETE")
	api.HandleFunc("/items/{id}", itemHandler.GetItemById).Methods("GET")
	api.HandleFunc("/items", itemHandler.ListItems).Methods("GET")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down server...")
		cancel()
	}()

	log.Println("Server running on port 8080")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error running server: %v", err)
		}
	}()
	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exiting")

}
