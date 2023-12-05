package models

import "gorm.io/gorm"

type Notification struct {
	gorm.Model
	Id          int    `json:"id" gorm:"primary_key"`
	Description string `json:"description"`
	Type        string `json:"Type"`
}
