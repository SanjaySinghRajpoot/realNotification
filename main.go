package main

import (
	"github.com/SanjaySinghRajpoot/realNotification/config"
	"github.com/SanjaySinghRajpoot/realNotification/routes"
	"github.com/SanjaySinghRajpoot/realNotification/utils"
	"github.com/gin-gonic/gin"
	"gopkg.in/robfig/cron.v2"
)

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

	routes.UserRoute(r)

	r.Run(":8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
