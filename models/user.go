package models

import "gorm.io/gorm"

type Notification struct {
	gorm.Model
	Id          int    `json:"id" gorm:"primary_key"`
	UserID      int    `json:"user_id"`
	Description string `json:"description" gorm:"type:varchar(1000)"`
	Type        string `json:"Type"`
	State       bool   `json:"state" gorm:"default:false"`
}

type NotificationPayload struct {
	// Type -> email, sms, inapp
	UserID      []int  `json:"user_id"`
	Type        string `json:"type"  validate:"required,oneof=sms email inapp"`
	Description string `json:"description"`
}

type NotificationValue struct {
	NotificationID int    `json:"id"`
	UserID         int    `json:"user_id"`
	Description    string `json:"description"`
}

type ServicePayload struct {
	Notification_id int    `json:"notification_id"`
	Message         string `json:"message"`
	UserID          int    `json:"user_id"`
}
