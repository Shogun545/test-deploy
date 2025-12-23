package entity

import "gorm.io/gorm"

type Major struct {
	gorm.Model
	Major string `json:"major"  binding:"required"`
}