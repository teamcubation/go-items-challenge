package main

import (
	"context"
	"fmt"
	"github.com/teamcubation/go-items-challenge/internal/adapters/http/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/teamcubation/go-items-challenge/internal/adapters/client"
	httphdl "github.com/teamcubation/go-items-challenge/internal/adapters/http"
	"github.com/teamcubation/go-items-challenge/internal/adapters/repository"
	"github.com/teamcubation/go-items-challenge/internal/application"
	"github.com/teamcubation/go-items-challenge/internal/domain/item"
	"github.com/teamcubation/go-items-challenge/internal/domain/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func runMigrations(db *gorm.DB) {
	err := db.AutoMigrate(&user.User{}, &item.Item{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}

func main() {
	err := godotenv.Load("/app/.env")
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
	categoryClient := client.NewCategoryClient("http://mockapi:8000")
	itemSrv := application.NewItemService(itemRepo, categoryClient)
	itemHandler := httphdl.NewItemHandler(itemSrv)
	exportHandler := httphdl.NewExportHandler(itemSrv)

	r := mux.NewRouter()

	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)

	api.HandleFunc("/items", itemHandler.CreateItem).Methods("POST")
	api.HandleFunc("/items/{id}", itemHandler.UpdateItem).Methods("PUT")
	api.HandleFunc("/items/{id}", itemHandler.DeleteItem).Methods("DELETE")
	api.HandleFunc("/items/{id}", itemHandler.GetItemByID).Methods("GET")
	api.HandleFunc("/items", itemHandler.ListItems).Methods("GET")
	api.HandleFunc("/items/export", exportHandler.Export).Methods("GET")

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
		log.Printf("Server forced to shutdown: %v", err)
		return
	}
	log.Println("Server exiting")
}
