package models

import "gorm.io/gorm"

type ProfessionGroup struct {
	gorm.Model
	Name   string  `json:"name" gorm:"unique;not null"`
	Titles []Title `json:"titles"`
}
