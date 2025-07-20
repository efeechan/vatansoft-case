package models

import (
	"gorm.io/gorm"
)

type Doctor struct {
	gorm.Model
	Name         string `json:"name" gorm:"not null"`
	Email        string `json:"email" gorm:"unique;not null"`
	Password     string `gorm:"not null"`
	HospitalID   uint   `json:"hospital_id"`
	Hospital     Hospital
	DepartmentID uint `json:"department_id"`
	Department   Department
}
