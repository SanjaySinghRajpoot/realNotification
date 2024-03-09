package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/SanjaySinghRajpoot/realNotification/config"
	"github.com/SanjaySinghRajpoot/realNotification/models"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var RedisClient *redis.Client
var KafkaProducer *kafka.Producer

func HomepageHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to Real notification"})
}

func SendNotification(Topic string, NotificationData models.NotificationValue, producer *kafka.Producer) (string, error) {

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

func CheckForNotificationState() {

	var allNotification []models.Notification

	res := config.DB.Debug().Where("state = false AND description !='' AND type !=''").Find(&allNotification)

	if res != nil {
		if res.RowsAffected == 0 {
			fmt.Println("No Notification found with False State")
			return
		}

		if res.Error != nil && res.RowsAffected != 0 {
			log := fmt.Sprintf("Error unable to fetch the data from DB: %s", res.Error)
			fmt.Println(log)
			return
		}
	}

	for _, notification := range allNotification {
		// Produce messages to the topic
		topic := notification.Type

		makeNotify := models.NotificationValue{
			NotificationID: notification.Id,
			UserID:         notification.UserID,
			Description:    notification.Description,
		}

		msg, err := SendNotification(topic, makeNotify, KafkaProducer)
		if err != nil {
			log := fmt.Sprintf("Unable to send Failed Notification from CRON %s", msg)
			fmt.Println(log)
			return
		}
	}
}

func SetUpRedis(password string) *redis.Client {

	return redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: password,
		DB:       0,
	})

}

func SetRedisData(UserID int, Description string, Type string) (string, error) {

	ctx := context.Background()

	// Unique key using userID and Description
	stringID := fmt.Sprintf("%d+%s", UserID, Description)

	err := RedisClient.Set(ctx, stringID, true, 3*time.Hour).Err()

	if err != nil {
		return "Something went wrong", err
	}

	return "", nil
}

func GetRedisData(UserID int, Description string) (bool, error) {

	ctx := context.Background()

	// Making unique key
	stringID := fmt.Sprintf("%d+%s", UserID, Description)

	check, err := RedisClient.Get(ctx, stringID).Bool()
	if err != nil {
		return false, err
	}

	return check, nil
}

func SetIPAddress(IPaddr string, Count int) (string, error) {

	ctx := context.Background()

	err := RedisClient.Set(ctx, IPaddr, Count, 30*time.Minute).Err()

	if err != nil {
		return "Something went wrong", err
	}

	return "", nil
}

func GetIPAddress(IPaddr string) (int, error) {

	ctx := context.Background()

	cnt, err := RedisClient.Get(ctx, IPaddr).Int()

	if err != nil {
		return -1, err
	}

	return cnt, nil
}

func InitializeProducer() (*kafka.Producer, error) {

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"client.id":         "test",
		"acks":              "all"})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		return nil, err
	}

	return producer, nil
}

func GetDB(userId int) *gorm.DB {

	if userId%2 == 0 {
		return config.DB
	}

	return config.DB1
}
