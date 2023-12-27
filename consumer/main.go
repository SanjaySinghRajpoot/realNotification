package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
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
	topics := []string{"sms", "email", "inapp"}
	consumer.SubscribeTopics(topics, nil)

	// Handle messages and shutdown signals
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
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
				fmt.Printf("Received message on topic %s: %s\n", *e.TopicPartition.Topic, string(e.Value))
				messageType := e.TopicPartition.Topic
				switch *messageType {
				case "sms":
					handleSMS(e)
				case "email":
					handleEmail(e)
				case "inapp":
					handleInapp(e)
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

func handleSMS(e *kafka.Message) {
	fmt.Printf("handleSMS %v\n", e)
}

func handleEmail(e *kafka.Message) {
	fmt.Printf("handleEmail %v\n", e)
}

func handleInapp(e *kafka.Message) {
	fmt.Printf("handleInapp %v\n", e)
}
