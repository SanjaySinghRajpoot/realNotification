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

	// 3 Different Consumers for three different functions

	// Create consumer for SMS service
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	// Subscribe to a topics
	topics := []string{"sms"}
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
			ev := consumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("Received message on topic %s: %s\n", *e.TopicPartition.Topic, string(e.Value))

			case kafka.Error:
				fmt.Fprintf(os.Stderr, "Error: %v\n", e)
				run = false

			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}

	// -------------------------------------------

	// Create consumer for Email
	consumerEmail, err := kafka.NewConsumer(config)
	if err != nil {
		panic(err)
	}
	defer consumerEmail.Close()

	// Subscribe to a topics
	topicsEmail := []string{"email"}
	consumerEmail.SubscribeTopics(topicsEmail, nil)

	// Handle messages and shutdown signals
	sigchanEmail := make(chan os.Signal, 1)
	signal.Notify(sigchanEmail, syscall.SIGINT, syscall.SIGTERM)

	runEmail := true
	for runEmail {
		select {
		case sig := <-sigchanEmail:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false

		default:
			ev := consumerEmail.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("Received message on topic %s: %s\n", *e.TopicPartition.Topic, string(e.Value))

			case kafka.Error:
				fmt.Fprintf(os.Stderr, "Error: %v\n", e)
				run = false

			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}
}
