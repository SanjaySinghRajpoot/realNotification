package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/SanjaySinghRajpoot/realNotification/config"
	"github.com/SanjaySinghRajpoot/realNotification/models"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func HomepageHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to Real notification"})
}

func CheckForNotificationState() {

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

	for _, notifi := range allNotification {
		// Produce messages to the topic
		topic := notifi.Type

		makeNotify := models.NotificationValue{
			NotificationID: notifi.Id,
			UserID:         notifi.UserID,
			Description:    notifi.Description,
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
			fmt.Printf("Failed to produce message 4: %v\n", err)
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

func SetUpRedis(password string) *redis.Client {

	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: password,
		DB:       0,
	})

}

type RedisNotification struct {
	Description string `redis:"description"`
	Type        string `redis:"type"`
}

func SetRedisData(UserID int, Description string, Type string) (string, error) {

	ctx := context.Background()

	stringID := fmt.Sprintf("%d", UserID)

	err := RedisClient.HSet(ctx, stringID, RedisNotification{Description, Type}).Err()

	if err != nil {
		return "Something went wrong", err
	}

	return "", nil
}

func GetRedisData(UserID int) (string, error) {

	ctx := context.Background()

	stringID := fmt.Sprintf("%d", UserID)

	var redisObj RedisNotification

	err := RedisClient.HGetAll(ctx, stringID).Scan(&redisObj)
	if err != nil {
		return "Something went wrong", err
	}

	return redisObj.Description, nil
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
