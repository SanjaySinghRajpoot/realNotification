package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomepageHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to Real notification"})
}

func CheckForNotificationState() {

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

	// var allNotification []models.Notification

	// res := config.DB.Where("state = false AND description !='' AND type !=''").Find(&allNotification)

	// if res != nil {
	// 	if res.RowsAffected == 0 {
	// 		fmt.Println("No Notification found with False State")
	// 		return
	// 	}

	// 	if res.Error != nil && res.RowsAffected != 0 {
	// 		log := fmt.Sprintf("Error unable to fetch the data from DB: %s", res.Error)
	// 		fmt.Println(log)
	// 		return
	// 	}
	// }

	// for _, notifi := range allNotification {
	// 	// Produce messages to the topic
	// 	topic := notifi.Type

	// 	makeNotify := models.NotificationValue{
	// 		ID:          notifi.Id,
	// 		Description: notifi.Description,
	// 	}

	// 	// Convert struct to bytes
	// 	notifyBytes, err := json.Marshal(makeNotify)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	deliveryChan := make(chan kafka.Event)
	// 	err = producer.Produce(&kafka.Message{
	// 		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
	// 		Value:          notifyBytes,
	// 	}, deliveryChan)

	// 	if err != nil {
	// 		fmt.Printf("Failed to produce message 4: %v\n", err)
	// 	} else {
	// 		// Wait for delivery report
	// 		e := <-deliveryChan
	// 		m := e.(*kafka.Message)
	// 		if m.TopicPartition.Error != nil {
	// 			fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	// 		} else {
	// 			fmt.Printf("Delivered message to topic %s [%d] at offset %v\n", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	// 		}
	// 	}
	// }

}
