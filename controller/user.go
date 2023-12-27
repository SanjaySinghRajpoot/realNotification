package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/SanjaySinghRajpoot/realNotification/config"
	"github.com/SanjaySinghRajpoot/realNotification/models"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
)

// Notification service
func Notification(ctx *gin.Context) {

	// Create producer
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"client.id":         "test",
		"acks":              "all"})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
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

	fmt.Println("----------------------------------------------")
	fmt.Println(notifyObj.Id)
	fmt.Println("----------------------------------------------")

	if res.Error != nil {
		fmt.Printf("Failed to create block: %v", res.Error)
	}

	if notificationPayload.Type == "sms" {

		// Produce messages to the topic
		topic := "sms"

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
			fmt.Printf("Failed to produce message: %v\n", err)
		} else {
			// Wait for delivery report
			e := <-deliveryChan
			m := e.(*kafka.Message)
			if m.TopicPartition.Error != nil {
				fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
			} else {
				fmt.Printf("Delivered message to topic %s [%d] at offset %v\n", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
			}
		}

	} else if notificationPayload.Type == "email" {
		// Produce messages to the topic
		topic := "email"

		deliveryChan := make(chan kafka.Event)
		err = producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(notificationPayload.Description),
		}, deliveryChan)

		if err != nil {
			fmt.Printf("Failed to produce message: %v\n", err)
		} else {
			// Wait for delivery report
			e := <-deliveryChan
			m := e.(*kafka.Message)
			if m.TopicPartition.Error != nil {
				fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
			} else {
				fmt.Printf("Delivered message to topic %s [%d] at offset %v\n", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
			}
		}
	} else if notificationPayload.Type == "inapp" {
		// Produce messages to the topic
		topic := "inapp"

		deliveryChan := make(chan kafka.Event)
		err = producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(notificationPayload.Description),
		}, deliveryChan)

		if err != nil {
			fmt.Printf("Failed to produce message: %v\n", err)
		} else {
			// Wait for delivery report
			e := <-deliveryChan
			m := e.(*kafka.Message)
			if m.TopicPartition.Error != nil {
				fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
			} else {
				fmt.Printf("Delivered message to topic %s [%d] at offset %v\n", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
			}
		}

	}

}
