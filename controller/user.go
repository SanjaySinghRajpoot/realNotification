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

	// Create producer
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"client.id":         "test",
		"acks":              "all"})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer producer.Close()

	var notificationPayload models.NotificationPayload

	ctx.BindJSON(&notificationPayload)

	// save the notification in the DB
	notifyObj := models.Notification{
		Type:        notificationPayload.Type,
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

	if notificationPayload.Type == utils.SMS {

		// Produce messages to the topic
		topic := utils.SMS

		makeNotify := models.NotificationValue{
			ID:          notifyObj.Id,
			Description: notificationPayload.Description,
		}

		// Convert struct to bytes
		notifyBytes, err := json.Marshal(makeNotify)
		if err != nil {
			log.Fatal(err)
		}

		deliveryChan := make(chan kafka.Event)
		err = producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          notifyBytes,
		}, deliveryChan)

		if err != nil {
			fmt.Printf("Failed to produce message 1: %v\n", err)

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return

		} else {
			// Wait for delivery report
			e := <-deliveryChan
			m := e.(*kafka.Message)
			if m.TopicPartition.Error != nil {
				fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)

				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": m.TopicPartition.Error.Error(),
				})
				return

			} else {
				fmt.Printf("Delivered message to topic %s [%d] at offset %v\n", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
			}
		}

	} else if notificationPayload.Type == utils.EMAIL {
		// Produce messages to the topic
		topic := utils.EMAIL

		deliveryChan := make(chan kafka.Event)
		err = producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(notificationPayload.Description),
		}, deliveryChan)

		if err != nil {
			fmt.Printf("Failed to produce message 2: %v\n", err)

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		} else {
			// Wait for delivery report
			e := <-deliveryChan
			m := e.(*kafka.Message)
			if m.TopicPartition.Error != nil {
				fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": m.TopicPartition.Error.Error(),
				})
				return

			} else {
				fmt.Printf("Delivered message to topic %s [%d] at offset %v\n", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
			}
		}
	} else if notificationPayload.Type == utils.PUSH {
		// Produce messages to the topic
		topic := utils.PUSH

		deliveryChan := make(chan kafka.Event)
		err = producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(notificationPayload.Description),
		}, deliveryChan)

		if err != nil {
			fmt.Printf("Failed to produce message 3: %v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		} else {
			// Wait for delivery report
			e := <-deliveryChan
			m := e.(*kafka.Message)
			if m.TopicPartition.Error != nil {
				fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": m.TopicPartition.Error.Error(),
				})
				return

			} else {
				fmt.Printf("Delivered message to topic %s [%d] at offset %v\n", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
			}
		}

	}

	ctx.JSON(http.StatusOK, "Notification was delivered successfully")
	return
}
