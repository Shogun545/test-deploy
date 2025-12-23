package repository

import (
	"backend/internal/app/entity"
	"gorm.io/gorm"
)

type ProfileRepository struct {
	DB *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{DB: db}
}

// ดึง user พร้อม relations ที่จำเป็นทั้งหมด (Prefix, Major, StudentProfile + AdvisorProfile.User)
func (r *ProfileRepository) FindUserWithProfileBySutID(sutID string) (*entity.User, error) {
	var user entity.User
	err := r.DB.
		Preload("Prefix").
		Preload("Major").
		Preload("StudentProfile.AdvisorProfile.User").
		Where("sut_id = ?", sutID).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *ProfileRepository) FindPrefixByName(name string) (*entity.Prefix, error) {
	var prefix entity.Prefix
	if err := r.DB.Where("prefix = ?", name).First(&prefix).Error; err != nil {
		return nil, err
	}
	return &prefix, nil
}

// ดึง StudentAcademicRecord ล่าสุดของ StudentProfile (อาจคืน ErrRecordNotFound)
func (r *ProfileRepository) FindLatestAcademicRecord(studentProfileID uint) (*entity.StudentAcademicRecord, error) {
	var rec entity.StudentAcademicRecord
	err := r.DB.
		Where("student_profile_id = ?", studentProfileID).
		Order("academic_year desc").
		Order("semester desc").
		First(&rec).Error
	if err != nil {
		return nil, err
	}
	return &rec, nil
}

// บันทึก User และ StudentProfile ใน transaction (caller ให้ DB ผ่าน)
func (r *ProfileRepository) SaveUser(tx *gorm.DB, user *entity.User) error {
	return tx.Save(user).Error
}

func (r *ProfileRepository) SaveStudentProfile(tx *gorm.DB, sp *entity.StudentProfile) error {
	return tx.Save(sp).Error
}
