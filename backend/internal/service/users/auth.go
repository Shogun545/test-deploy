package service

import (
	"backend/config"
	"backend/internal/app/entity"
	"backend/internal/app/repository"
	"errors"
	"time"
)

type AuthService struct {
	UserRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{UserRepo: userRepo}
}

func (s *AuthService) Login(sutID, password string) (*entity.User, error) {
	if sutID == "" {
		return nil, errors.New("sut_id is required")
	}
	user, err := s.UserRepo.FindBySutID(sutID)
	if err != nil {
		return nil, errors.New("invalid sut_id or password")
	}
	// ใช้ config.CheckPasswordHash(plain, hash)
	if !config.CheckPasswordHash([]byte(password), []byte(user.PasswordHash)) {
		return nil, errors.New("invalid sut_id or password")
	}
	now := time.Now()
	user.LastLogin = &now

	if err := s.UserRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
