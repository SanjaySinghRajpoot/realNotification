package consumer

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SanjaySinghRajpoot/realNotification/config"
	"github.com/SanjaySinghRajpoot/realNotification/models"
	"github.com/SanjaySinghRajpoot/realNotification/utils"
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
					handleSMS(e, notifObj.ID)
				case utils.EMAIL:
					handleEmail(e, notifObj.ID)
				case utils.PUSH:
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
