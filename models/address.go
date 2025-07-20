package models

import "gorm.io/gorm"

type Address struct {
	gorm.Model
	Street     string `json:"street"`
	City       string `json:"city"`
	ProvinceID uint   `json:"province_id"`
	DistrictID uint   `json:"district_id"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}
