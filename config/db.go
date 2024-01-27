package config

import (
	"log"
	"os"

	"github.com/SanjaySinghRajpoot/realNotification/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func Connect() {
	// dsn := os.Getenv("DATABASE_URL")
	// dsn := os.Getenv("host=localhost user=postgres password=postgres dbname=postgres sslmode=disable")
	dsn := "host=localhost user=postgres password=postgres dbname=postgres sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.Notification{})

	DB = db
}
