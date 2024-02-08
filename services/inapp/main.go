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

type Payload struct {
	Notification_id int    `json:"notification_id"`
	Message         string `json:"message"`
	UserID          int    `json:"user_id"`
}

var DB *gorm.DB
var DB1 *gorm.DB

func Connect() {
	// dsn := os.Getenv("DATABASE_URL")
	// dsn := os.Getenv("host=localhost user=postgres password=postgres dbname=postgres sslmode=disable")
	dsn := "host=localhost user=postgres password=postgres dbname=postgres sslmode=disable"

	dsn1 := "host=localhost port=5433 user=postgres password=postgres dbname=postgres sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db1, err := gorm.Open(postgres.Open(dsn1), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	DB = db
	DB1 = db1
}

func GetDB(userId int) *gorm.DB {

	if userId%2 == 0 {
		return DB
	}

	return DB1
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
	useDB := GetDB(int(payload.UserID))
	res := useDB.Model(&updateNotification).Where("id = ? AND user_id = ?", payload.Notification_id, payload.UserID).Update("state", true)

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
