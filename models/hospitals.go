package models

import "gorm.io/gorm"

type Hospital struct {
	gorm.Model
	Name      string  `json:"name"`
	TaxNumber string  `json:"tax_number"`
	Email     string  `json:"email"`
	Phone     string  `json:"phone"`
	AddressID uint    `json:"address_id"`
	Address   Address `gorm:"foreignKey:AddressID" json:"address"`
	Users     []User  `json:"users,omitempty"`
}
