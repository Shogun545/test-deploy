package repository

import (
	"backend/internal/app/entity"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) FindBySutID(sutID string) (*entity.User, error) {
	var user entity.User
	err := r.DB.        
		Preload("Role").
		Where("sut_id = ?", sutID).
		First(&user).Error
	return &user, err
}
func (r *UserRepository) Update(user *entity.User) error {
	return r.DB.Save(user).Error
}

