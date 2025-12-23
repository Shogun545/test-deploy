package entity

import "gorm.io/gorm"

type Prefix struct {
	gorm.Model
	Prefix string `json:"prefix"  binding:"required"`
}