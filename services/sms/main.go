package main

import (
	"fmt"
	"net/http"

	"github.com/SanjaySinghRajpoot/realNotification/config"
	"github.com/SanjaySinghRajpoot/realNotification/models"
	"github.com/gin-gonic/gin"
)

func HomepageHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to Real SMS notification"})
}

func SMSService(ctx *gin.Context) {

	var payload models.SMSpayload

	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("%v", err)})
		return
	}

	// Update the state of the notification based on service used
	var updateNotification models.Notification

	// we need to get the ID from the payload
	res := config.DB.Model(&updateNotification).Where("id = ?", payload.Notification_id).Update("state", true)

	if res.Error != nil {
		fmt.Printf("Failed to update the Notification: %v", res.Error)
	}
}

func main() {

	// Connect to the database
	config.Connect()

	// Gin router
	r := gin.Default()

	r.GET("/", HomepageHandler)

	r.POST("/sms", SMSService)

	r.Run(":8082") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
