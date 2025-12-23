package entity

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Role string `json:"role"  binding:"required"`
}