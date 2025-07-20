package models

import "gorm.io/gorm"

type Department struct {
	gorm.Model
	Name             string `gorm:"not null"`
	HospitalID       uint
	Hospital         Hospital `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	DepartmentTypeID uint
	DepartmentType   DepartmentType `gorm:"foreignKey:DepartmentTypeID"`
}

type DepartmentType struct {
	gorm.Model
	Name string `gorm:"unique"`
}
