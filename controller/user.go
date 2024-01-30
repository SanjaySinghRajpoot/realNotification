package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/SanjaySinghRajpoot/realNotification/config"
	"github.com/SanjaySinghRajpoot/realNotification/models"
	"github.com/SanjaySinghRajpoot/realNotification/utils"
	"github.com/confluentinc/confluent-kafka-go/kafka"
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

		// save the notification in the DB
		notifyObj := models.Notification{
			Type:        notificationPayload.Type,
			UserID:      userID,
			Description: notificationPayload.Description,
			State:       false,
		}

		res := config.DB.Create(&notifyObj)

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
			msg, err := SendNotification(utils.SMS, notificationKafkaObj, utils.KafkaProducer)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
			}

			ctx.JSON(http.StatusOK, msg)
		} else if notificationPayload.Type == utils.EMAIL {
			msg, err := SendNotification(utils.EMAIL, notificationKafkaObj, utils.KafkaProducer)

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
			}

			ctx.JSON(http.StatusOK, msg)
		} else if notificationPayload.Type == utils.PUSH {
			msg, err := SendNotification(utils.PUSH, notificationKafkaObj, utils.KafkaProducer)

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

func SendNotification(Topic string, NotificationData models.NotificationValue, producer *kafka.Producer) (string, error) {

	defer producer.Close()

	// Convert struct to bytes
	notifyBytes, err := json.Marshal(NotificationData)
	if err != nil {
		log.Fatal(err)
	}

	deliveryChan := make(chan kafka.Event)
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &Topic, Partition: kafka.PartitionAny},
		Value:          notifyBytes,
	}, deliveryChan)

	if err != nil {
		msg := fmt.Sprintf("Failed to produce message 1: %v\n", err)

		return msg, err

	} else {

		// Wait for delivery report
		e := <-deliveryChan
		m := e.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			msg := fmt.Sprintf("Delivery failed: %v\n", m.TopicPartition.Error)

			return msg, m.TopicPartition.Error

		} else {
			fmt.Printf("Delivered message to topic %s [%d] at offset %v\n", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
		}
	}

	return "Message Delivered Successfully", nil
}
