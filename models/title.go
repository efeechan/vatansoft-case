package models

import "gorm.io/gorm"

type Title struct {
	gorm.Model
	Name              string `json:"name" gorm:"not null"`
	ProfessionGroupID uint   `json:"profession_group_id"`
}
