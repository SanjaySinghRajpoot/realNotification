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
				case utils.SMS:
					handleSMS(e, notifObj.NotificationID, notifObj.UserID)
				case utils.EMAIL:
					handleEmail(e, notifObj.NotificationID, notifObj.UserID)
				case utils.PUSH:
					handleInapp(e, notifObj.NotificationID, notifObj.UserID)
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

func handleSMS(e *kafka.Message, notifID int, userID int) {

	url := "http://localhost:8082/sms"

	payload := models.ServicePayload{
		Notification_id: notifID,
		UserID:          userID,
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

func handleEmail(e *kafka.Message, notifID int, userID int) {

	url := "http://localhost:8083/mail"

	payload := models.ServicePayload{
		Notification_id: notifID,
		UserID:          userID,
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

	fmt.Printf("handleEmail %v\n", e)
}

func handleInapp(e *kafka.Message, notifID int, userID int) {
	url := "http://localhost:8084/inapp"

	payload := models.ServicePayload{
		Notification_id: notifID,
		UserID:          userID,
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

	fmt.Printf("handleInapp %v\n", e)
}
