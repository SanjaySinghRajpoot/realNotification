package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	Id          int    `json:"id" gorm:"primary_key"`
	Description string `json:"description"`
	Type        string `json:"Type"`
	State       bool   `json:"state" gorm:"default:false"`
}

type Payload struct {
	Notification_id int    `json:"notification_id"`
	Message         string `json:"message"`
}

var DB *gorm.DB

func Connect() {
	// dsn := "host=localhost user=postgres password=postgres dbname=notification port=5432 sslmode=disable"
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	DB = db
}

func HomepageHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to Real Inapp notification"})
}

func InappService(ctx *gin.Context) {

	var payload Payload

	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("%v", err)})
		return
	}

	fmt.Println(payload.Message)

	// Update the state of the notification based on service used
	var updateNotification Notification

	// we need to get the ID from the payload
	res := DB.Model(&updateNotification).Where("id = ?", payload.Notification_id).Update("state", true)

	if res.Error != nil {
		fmt.Printf("Failed to update the Notification: %v", res.Error)
	}
}

func main() {

	// Connect to the database
	Connect()

	// Gin router
	r := gin.Default()

	r.GET("/", HomepageHandler)

	r.POST("/inapp", InappService)

	r.Run(":8084")
}
