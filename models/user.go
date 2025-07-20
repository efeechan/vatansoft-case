package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Email      string `json:"email" gorm:"unique"`
	Password   string `json:"-"`
	Phone      string `json:"phone" gorm:"unique"`
	TCKN       string `json:"tckn" gorm:"unique"`
	Role       string `json:"role"`
	HospitalID uint   `json:"hospital_id"`
	Hospital   Hospital

	ProfessionGroupID uint             `json:"profession_group_id"`
	ProfessionGroup   *ProfessionGroup `json:"profession_group,omitempty" gorm:"foreignKey:ProfessionGroupID"`

	TitleID uint   `json:"title_id"`
	Title   *Title `json:"title,omitempty" gorm:"foreignKey:TitleID"`
}
