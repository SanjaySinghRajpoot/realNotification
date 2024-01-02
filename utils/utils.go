package utils

import (
	"fmt"

	"github.com/SanjaySinghRajpoot/realNotification/config"
	"github.com/SanjaySinghRajpoot/realNotification/models"
)

// // Create producer
// producer, err := kafka.NewProducer(&kafka.ConfigMap{
// 	"bootstrap.servers": "localhost:9092",
// 	"client.id":         "test",
// 	"acks":              "all"})

// if err != nil {
// 	fmt.Printf("Failed to create producer: %s\n", err)
// 	os.Exit(1)
// }
// defer producer.Close()

func CheckForNotificationState() {

	var allNotification []models.Notification

	res := config.DB.Where("state = false").Find(&allNotification)

	if res != nil {
		fmt.Println("Error unable to fetch the data from DB")
	}

	fmt.Println("cron is working")

	// Produce messages to the topic
	topic := "sms"

	fmt.Println(topic)

	// makeNotify := models.NotificationValue{
	// 	ID:          notifyObj.Id,
	// 	Description: notificationPayload.Description,
	// }

	// // Convert struct to bytes
	// notifyBytes, err := json.Marshal(makeNotify)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// deliveryChan := make(chan kafka.Event)
	// err = producer.Produce(&kafka.Message{
	// 	TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
	// 	Value:          notifyBytes,
	// }, deliveryChan)

	// if err != nil {
	// 	fmt.Printf("Failed to produce message: %v\n", err)
	// } else {
	// 	// Wait for delivery report
	// 	e := <-deliveryChan
	// 	m := e.(*kafka.Message)
	// 	if m.TopicPartition.Error != nil {
	// 		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	// 	} else {
	// 		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	// 	}
	// }
}
