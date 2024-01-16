package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/SanjaySinghRajpoot/realNotification/config"
	"github.com/SanjaySinghRajpoot/realNotification/models"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {

	// Connect to the postgres database
	config.Connect()

	// Set up configuration
	config := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092", // Replace with your Kafka broker address
		"group.id":          "my-group",
		"auto.offset.reset": "earliest",
	}

	// Create consumer for all three types
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	// Subscribe to a topics
	// topics := []string{"sms", "email", "inapp"}
	topics := []string{"sms"}
	consumer.SubscribeTopics(topics, nil)

	// Handle messages and shutdown signals
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	var notifObj models.NotificationValue
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false

		default:
			// time out of 100 millisecond
			ev := consumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:

				err = json.Unmarshal(e.Value, &notifObj)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf("Received message on topic %s: %s\n", *e.TopicPartition.Topic, notifObj.Description)
				messageType := e.TopicPartition.Topic
				switch *messageType {
				case "sms":
					handleSMS(e, notifObj.ID)
				case "email":
					handleEmail(e, notifObj.ID)
				case "inapp":
					handleInapp(e, notifObj.ID)
				}

			case kafka.Error:
				fmt.Fprintf(os.Stderr, "Error: %v\n", e)
				run = false

			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}
}

func handleSMS(e *kafka.Message, notifID int) {

	// // Update the state of the notification based on service used
	// var updateNotification models.Notification

	// res := config.DB.Model(&updateNotification).Where("id = ?", notifID).Update("state", true)

	// if res.Error != nil {
	// 	fmt.Printf("Failed to update the Notification: %v", res.Error)
	// }

	// need to set the status of the notification to true in the DB

	url := "http://localhost:8082/sms"

	payload := models.SMSpayload{
		Notification_id: notifID,
		Message:         string(e.Value),
	}

	jsonStr, err := json.Marshal(&payload)
	if err != nil {
		fmt.Println("Error while Marshalling:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

}

func handleEmail(e *kafka.Message, notifID int) {

	// Update the state of the notification based on service used
	var updateNotification models.Notification

	res := config.DB.Model(&updateNotification).Where("id = ?", notifID).Update("state", true)

	if res.Error != nil {
		fmt.Printf("Failed to update the Notification: %v", res.Error)
	}

	fmt.Printf("handleEmail %v\n", e)
}

func handleInapp(e *kafka.Message, notifID int) {

	// Update the state of the notification based on service used

	var updateNotification models.Notification

	res := config.DB.Model(&updateNotification).Where("id = ?", notifID).Update("state", true)

	if res.Error != nil {
		fmt.Printf("Failed to update the Notification: %v", res.Error)
	}

	fmt.Printf("handleInapp %v\n", e)
}
