package repository

import (
	"backend/internal/app/entity"
	"gorm.io/gorm"
	"errors"
)

// interface
type AdvisorProfileRepository interface {
    FindBySutID(sutID string) (*entity.AdvisorProfile, error)
    Save(profile *entity.AdvisorProfile) error
	FindAdvisorWithStudentsBySutID(sutID string) (*entity.User, error)
    FindPrefixByName(name string) (*entity.Prefix, error)
}


// struct
type advisorProfileRepository struct {
    db *gorm.DB
}

// constructor
func NewAdvisorProfileRepository(db *gorm.DB) AdvisorProfileRepository {
    return &advisorProfileRepository{db: db}
}

// implement interface สำหรับ profile
func (r *advisorProfileRepository) FindBySutID(sutID string) (*entity.AdvisorProfile, error) {
    var profile entity.AdvisorProfile
    // join กับ users เพื่อ filter ด้วย sut_id
    err := r.db.
        Preload("User").
        Preload("User.Major").
        Preload("User.Prefix").
        Joins("JOIN users ON users.id = advisor_profiles.user_id").
        Where("users.sut_id = ?", sutID).
        First(&profile).Error
    if err != nil {
        return nil, errors.New("record not found")
    }
    return &profile, nil
}


func (r *advisorProfileRepository) Save(profile *entity.AdvisorProfile) error {
    return r.db.Save(profile).Error
}

// method เพิ่มสำหรับดึง advisor + students
func (r *advisorProfileRepository) FindAdvisorWithStudentsBySutID(sutID string) (*entity.User, error) {
    var advisor entity.User
    if err := r.db.
        Preload("AdvisorProfile.Students.User.Major").
        Preload("AdvisorProfile.Students.User.StudentProfile.StudentAcademicRecords", func(db *gorm.DB) *gorm.DB {
            return db.Order("academic_year desc, semester desc")
        }).
        Where("sut_id = ?", sutID). 
        First(&advisor).Error; err != nil {
        return nil, errors.New("record not found")
    }
    return &advisor, nil
}

func (r *advisorProfileRepository) FindAdvisorWithStudents(advisorSutID string) (*entity.User, error) {
	var advisor entity.User

	// preload ความสัมพันธ์ที่ต้องการให้ครบ (ปรับชื่อ relation ให้ตรง entity ของโปรเจกต์)
	err := r.db.
		Preload("AdvisorProfile.Students.User.StudentProfile.StudentAcademicRecords").
		Preload("AdvisorProfile.Students.User.Major").
		Where("sut_id = ?", advisorSutID).
		First(&advisor).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &advisor, nil
}


// แก้ไขบรรทัดนี้:
func (r *advisorProfileRepository) FindPrefixByName(name string) (*entity.Prefix, error) {
    var prefix entity.Prefix
    if err := r.db.Where("prefix = ?", name).First(&prefix).Error; err != nil {
        return nil, err
    }
    return &prefix, nil
}