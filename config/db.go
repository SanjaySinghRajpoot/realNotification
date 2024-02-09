package config

import (
	"os"

	"github.com/SanjaySinghRajpoot/realNotification/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var DB1 *gorm.DB

func Connect() {
	dsn := os.Getenv("DATABASE_URL")
	dsn1 := os.Getenv("DATABASE_URL_1")
	// dsn := os.Getenv("host=localhost user=postgres password=postgres dbname=postgres sslmode=disable")
	// dsn := "host=localhost user=postgres password=postgres dbname=postgres sslmode=disable"

	// dsn1 := "host=localhost port=5433 user=postgres password=postgres dbname=postgres sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db1, err := gorm.Open(postgres.Open(dsn1), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.Notification{})
	db1.AutoMigrate(&models.Notification{})

	DB = db
	DB1 = db1
}
