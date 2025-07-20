package models

import "gorm.io/gorm"

type City struct {
	gorm.Model
	Name      string     `json:"name" gorm:"unique;not null"`
	Districts []District `json:"districts"`
}
