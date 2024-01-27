package main

import (
	"fmt"
	"net/http"

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

type ServicePayload struct {
	Notification_id int    `json:"notification_id"`
	Message         string `json:"message"`
	UserID          int    `json:"user_id"`
}

var DB *gorm.DB

func Connect() {
	// dsn := os.Getenv("DATABASE_URL")
	dsn := "host=localhost user=postgres password=postgres dbname=postgres sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	DB = db
}

func HomepageHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to Real SMS notification"})
}

func SMSService(ctx *gin.Context) {

	var payload ServicePayload

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("%v", err)})
		return
	}

	fmt.Println(payload.Message)

	// Update the state of the notification based on service used
	var updateNotification Notification

	// we need to get the ID from the payload
	res := DB.Model(&updateNotification).Where("id = ? AND user_id = ?", payload.Notification_id, payload.UserID).Update("state", true)

	if res.Error != nil {
		fmt.Printf("Failed to update the Notification: %v", res.Error)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "SMS sent",
	})
	return
}

func main() {

	// Connect to the database
	Connect()

	// Gin router
	r := gin.Default()

	r.GET("/", HomepageHandler)

	r.POST("/sms", SMSService)

	r.Run(":8082") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
