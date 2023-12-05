package routes

import (
	"github.com/SanjaySinghRajpoot/realNotification/controller"
	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {

	// All user routes
	router.POST("/user/notification", controller.Notification)
}
