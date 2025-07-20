package models

import "gorm.io/gorm"

type District struct {
	gorm.Model
	Name   string `json:"name" gorm:"not null"`
	CityID uint   `json:"city_id"`
}
