package models

import "gorm.io/gorm"

type Notification struct {
	gorm.Model
	Id          int    `json:"id" gorm:"primary_key"`
	Description string `json:"description"`
	Type        string `json:"Type"`
	State       bool   `json:"state" gorm:"default:false"`
}

type NotificationPayload struct {
	// Type -> email, sms, inapp
	Type        string `json:"type"`
	Description string `json:"description"`
}

type NotificationValue struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type SMSpayload struct {
	Notification_id int    `json:"notification_id"`
	Message         string `json:"message"`
}
