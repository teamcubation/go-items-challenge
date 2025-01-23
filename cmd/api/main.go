package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/teamcubation/go-items-challenge/cmd/api/server"
	"github.com/teamcubation/go-items-challenge/internal/domain/item"
	"github.com/teamcubation/go-items-challenge/internal/domain/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	config := loadConfig()

	db := initDatabase(config)

	runMigrations(db)

	r := server.SetupRouter(db, config.ExternalAPIURL)

	srv := server.NewServer(r, config.ServerPort)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := srv.Start(ctx); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}
}

type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	ExternalAPIURL string
	ServerPort     string
}

func loadConfig() Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return Config{
		DBHost:         os.Getenv("DB_HOST"),
		DBPort:         os.Getenv("DB_PORT"),
		DBUser:         os.Getenv("DB_USER"),
		DBPassword:     os.Getenv("DB_PASSWORD"),
		DBName:         os.Getenv("DB_NAME"),
		ExternalAPIURL: os.Getenv("EXTERNAL_API_URL"),
		ServerPort:     ":8080",
	}
}

func initDatabase(config Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return db
}

func runMigrations(db *gorm.DB) {
	err := db.AutoMigrate(&user.User{}, &item.Item{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
