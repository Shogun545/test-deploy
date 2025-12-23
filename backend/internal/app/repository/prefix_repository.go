package repository

import (
	"backend/internal/app/entity"
	"gorm.io/gorm"
)

type PrefixRepository interface {
	FindAll() ([]entity.Prefix, error)
}

type prefixRepository struct {
	db *gorm.DB
}

func NewPrefixRepository(db *gorm.DB) PrefixRepository {
	return &prefixRepository{db: db}
}

func (r *prefixRepository) FindAll() ([]entity.Prefix, error) {
	var prefixes []entity.Prefix
	err := r.db.Order("prefix ASC").Find(&prefixes).Error
	return prefixes, err
}
