package main

import (
	"github.com/SanjaySinghRajpoot/realNotification/config"
	"github.com/SanjaySinghRajpoot/realNotification/routes"
	"github.com/SanjaySinghRajpoot/realNotification/utils"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gopkg.in/robfig/cron.v2"
)

// @title 	Real Notification Service
// @version	1.0
// @description A Notification Service in Go using Gin framework

// @host 	localhost:8080
// @BasePath /
func main() {

	// Connect to the database
	config.Connect()

	// start the cron job
	cronJob := cron.New()

	// Cron Job that will check for state=false for Notifications
	cronJob.AddFunc("@every 10s", func() {
		utils.CheckForNotificationState()
	})

	cronJob.Start()

	// Gin router
	r := gin.Default()

	r.GET("/", utils.HomepageHandler)

	// setting up the SWAAGGER URL
	url := ginSwagger.URL("http://localhost:8081/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	routes.UserRoute(r)

	r.Run(":8081")
}
