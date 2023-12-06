package controller

import (
	"fmt"

	"github.com/SanjaySinghRajpoot/realNotification/models"
	"github.com/SanjaySinghRajpoot/realNotification/utils"
	"github.com/gin-gonic/gin"
)

// Notification service
func Notification(ctx *gin.Context) {
	fmt.Println("working--- MNC")

	var notificationPayload models.NotificationPayload

	ctx.BindJSON(&notificationPayload)

	if notificationPayload.Type == "sms" {

		smsStatus, checkSMS := utils.SendSMSNotification(notificationPayload.Description)

		if !checkSMS {
			ctx.JSON(404, fmt.Sprintf("Something went wrong while sending an SMS - %v", smsStatus))
		}
	} else if notificationPayload.Type == "email" {
		emailStatus, checkEmail := utils.SendEmailNotification(notificationPayload.Description)

		if !checkEmail {
			ctx.JSON(404, fmt.Sprintf("Something went wrong while sending an Email - %v", emailStatus))
		}
	} else if notificationPayload.Type == "inapp" {
		emailStatus, checkEmail := utils.SendInappNotification(notificationPayload.Description)

		if !checkEmail {
			ctx.JSON(404, fmt.Sprintf("Something went wrong while sending an Inapp msg - %v", emailStatus))
		}
	}

}
