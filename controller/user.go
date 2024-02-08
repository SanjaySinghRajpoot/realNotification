package controller

import (
	"fmt"
	"net/http"

	"github.com/SanjaySinghRajpoot/realNotification/models"
	"github.com/SanjaySinghRajpoot/realNotification/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Endpoint to accept all the notifications
// @Description This endpoint will accept all the notifications from all the services
// @Accept  json
// @Produce  json
// @Success 200 {string} string  "ok"
// @Router /user/notification [POST]
func Notification(ctx *gin.Context) {

	var notificationPayload models.NotificationPayload

	ctx.BindJSON(&notificationPayload)

	for _, userID := range notificationPayload.UserID {

		// check first in cache
		check, err := utils.GetRedisData(userID, notificationPayload.Description)
		if err != nil {
			fmt.Printf("Failed to Get the Redis Cache, Setting the Cache: %s", err)
		}

		if check {
			ctx.JSON(http.StatusOK, "This notification is already sent in the last 24 hours")
			return
		}

		msg, error := utils.SetRedisData(userID, notificationPayload.Description, notificationPayload.Type)

		if error != nil {

			fmt.Printf("Failed to Set the Redis Cache: %s", msg)

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": error.Error(),
			})

			return
		}

		useDB := utils.GetDB(userID)

		// save the notification in the DB
		notifyObj := models.Notification{
			Type:        notificationPayload.Type,
			UserID:      userID,
			Description: notificationPayload.Description,
			State:       false,
		}

		res := useDB.Create(&notifyObj)

		if res.Error != nil {
			fmt.Printf("Failed to create block: %v", res.Error)

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": res.Error.Error(),
			})
			return
		}

		notificationKafkaObj := models.NotificationValue{
			NotificationID: notifyObj.Id,
			UserID:         userID,
			Description:    notificationPayload.Description,
		}

		if notificationPayload.Type == utils.SMS {
			msg, err := utils.SendNotification(utils.SMS, notificationKafkaObj, utils.KafkaProducer)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
			}

			ctx.JSON(http.StatusOK, msg)
		} else if notificationPayload.Type == utils.EMAIL {
			msg, err := utils.SendNotification(utils.EMAIL, notificationKafkaObj, utils.KafkaProducer)

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
			}

			ctx.JSON(http.StatusOK, msg)
		} else if notificationPayload.Type == utils.PUSH {
			msg, err := utils.SendNotification(utils.PUSH, notificationKafkaObj, utils.KafkaProducer)

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
			}

			ctx.JSON(http.StatusOK, msg)
		}

	}

	return
}
