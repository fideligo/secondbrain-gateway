package database

import (
	"fmt"
	"log"
	"os"

	"github.com/fideligo/secondbrain-gateway/internal/model"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Retrieve credentials from environment variables
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// Construct DSN
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", 
		host, user, password, dbName, port)
	
	// Open connection to PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	log.Println("Successfully connected to PostgreSQL database")

	// Automatically create/update tables based on models
	err = db.AutoMigrate(&model.Document{})
	if err != nil {
		log.Fatal("Database migration failed: ", err)
	}

	log.Println("Database migration completed successfully")

	return db
}