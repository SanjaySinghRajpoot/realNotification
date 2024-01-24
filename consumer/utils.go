package consumer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SanjaySinghRajpoot/realNotification/models"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func handleSMS(e *kafka.Message, notifID int) {

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

	url := "http://localhost:8083/mail"

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

	fmt.Printf("handleEmail %v\n", e)
}

func handleInapp(e *kafka.Message, notifID int) {

	url := "http://localhost:8084/inapp"

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

	fmt.Printf("handleInapp %v\n", e)
}
