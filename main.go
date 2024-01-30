package main

import (
	"fmt"
	"log"
	"os"

	"github.com/SanjaySinghRajpoot/realNotification/config"
	"github.com/SanjaySinghRajpoot/realNotification/routes"
	"github.com/SanjaySinghRajpoot/realNotification/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gopkg.in/robfig/cron.v2"
)

// @title 	Real Notification Service
// @version	1.0
// @description A Notification Service in Go using Gin framework
// @host 	localhost:8081
// @BasePath /
func main() {

	// Connect to the database
	config.Connect()

	// start the CRON JOB
	CRONjob()

	// Gin router
	r := gin.Default()

	// Home Page endpoint
	r.GET("/", utils.HomepageHandler)

	password := EnvVariable("PASSWORD")

	// Redis Cache Setup
	utils.RedisClient = utils.SetUpRedis(password)

	var err error
	utils.KafkaProducer, err = utils.InitializeProducer()

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err.Error())
		return
	}

	// adding rate limiter for all the routes
	// r.Use(middleware.RateLimiter)

	// User Router for Notification Route
	routes.UserRoute(r)

	// setting up the SWAAGGER URL
	url := ginSwagger.URL("http://localhost:8081/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	r.Run(":8081")
}

func CRONjob() {
	// start the cron job
	cronJob := cron.New()

	// Cron Job that will check for state=false for Notifications
	cronJob.AddFunc("@every 10s", func() {
		utils.CheckForNotificationState()
	})

	cronJob.Start()
}

// use godot package to load/read the .env file and
// return the value of the key
func EnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
