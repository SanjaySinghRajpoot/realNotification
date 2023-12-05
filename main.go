package main

import (
	"github.com/SanjaySinghRajpoot/realNotification/config"
	"github.com/SanjaySinghRajpoot/realNotification/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	config.Connect()

	routes.UserRoute(r)

	r.Run(":8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
