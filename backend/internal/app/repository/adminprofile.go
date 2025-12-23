package repository

import (
	"backend/internal/app/dto"
	"backend/internal/app/entity"
	"errors"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("repository record not found")

// AdminProfileRepository Interface กำหนดสัญญาของฟังก์ชันที่ Admin Service ต้องการใช้
type AdminProfileRepository interface {
	// สำหรับ Get/Update MyProfile ของ Admin
	FindBySutID(sutID string) (*entity.User, error)

	// สำหรับ Update MyProfile (Save/Update)
	UpdateUser(user *entity.User) error

	// สำหรับ GetAllUsers
	FindAllUsers(filters dto.UserListFilter) ([]*entity.User, error)

	// สำหรับ GetUserBySutID: ดึงรายละเอียดผู้ใช้งานคนใดคนหนึ่ง
	FindUserDetailBySutID(sutID string) (*entity.User, error)

	FindAllMajors() ([]*entity.Major, error)

	UpdateUserStatus(sutID string, newStatus string) error

	FindUserBySutID(sutID string) (*entity.User, error)
}

// AdminProfileRepositoryImpl คือโครงสร้างที่ implement Interface
type AdminProfileRepositoryImpl struct {
	db *gorm.DB
}

// NewAdminProfileRepositoryImpl สร้าง instance ของ Repository
func NewAdminProfileRepositoryImpl(db *gorm.DB) AdminProfileRepository {
	return &AdminProfileRepositoryImpl{db: db}
}

// FindBySutID ดึงข้อมูล User (Admin) พร้อม Preload Department และ Role
func (r *AdminProfileRepositoryImpl) FindBySutID(sutID string) (*entity.User, error) {
	var user entity.User
	result := r.db.
		Preload("Department").
		Preload("Role").
		Preload("Prefix").
		Preload("ManagedUsers").
		Preload("ManagedUsers.Role").
		Preload("ManagedUsers.Prefix").
		Where("sut_id = ?", sutID).
		First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

// UpdateUser บันทึกการเปลี่ยนแปลงข้อมูล User
func (r *AdminProfileRepositoryImpl) UpdateUser(user *entity.User) error {
	result := r.db.Save(user)
	return result.Error
}

// FindAllUsers ดึง User ทั้งหมด (เพื่อแสดงในหน้า Managed Users) พร้อม Filter
func (r *AdminProfileRepositoryImpl) FindAllUsers(filters dto.UserListFilter) ([]*entity.User, error) {
	var users []*entity.User
	query := r.db.
		Preload("Role").
		Preload("Department").
		Preload("Major").
		Preload("Prefix")

	// กรองตาม Role
	if filters.Role != "" {
		// Joins: เข้าร่วมตาราง Role เพื่อกรองตามชื่อ Role
		query = query.Joins("JOIN roles ON roles.id = users.role_id").Where("roles.role = ?", filters.Role)
	}

	// กรองตาม Status (Active/Inactive)
	switch filters.Status {
	case "active":
		query = query.Where("users.active = ?", true)
	case "inactive":
		query = query.Where("users.active = ?", false)
	case "":
		// ไม่กรอง แสดงทั้งหมด
	default:
		// status ไม่ถูกต้อง → คืนค่าว่าง
		return []*entity.User{}, nil
	}

	//  Logic การกรอง Major ที่ถูกต้อง (ใช้ Major ID)
	if filters.Major != "" {
		// 1. หา Major ID จากชื่อ Major ที่ส่งมา
		var major entity.Major
		// (3) ลบ .Debug() ออกจาก Query ค้นหา Major
		err := r.db.Where("major = ?", filters.Major).First(&major).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// ถ้าหา Major ไม่เจอ, คืนรายการว่าง
				return []*entity.User{}, nil
			}
			return nil, err
		}

		// 2. กรองผู้ใช้ (users) ตาม major_id ที่หามาได้
		query = query.Where("users.major_id = ?", major.ID)
	}

	// ค้นหา (Search)
	if filters.Search != "" {
		// ค้นหาใน SutID, FirstName, LastName, หรือ Email
		searchTerm := "%" + filters.Search + "%"
		query = query.Where(
			"users.sut_id LIKE ? OR users.first_name LIKE ? OR users.last_name LIKE ? OR users.email LIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	// 3. รัน Query
	result := query.Find(&users)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return []*entity.User{}, nil
		}
		return nil, result.Error
	}
	return users, nil
}

// FindUserDetailBySutID ดึง User รายคน พร้อมข้อมูลรายละเอียดทั้งหมด (Student, Advisor, Academic)
func (r *AdminProfileRepositoryImpl) FindUserDetailBySutID(sutID string) (*entity.User, error) {
	var user entity.User

	// Preload ข้อมูลรายละเอียดทั้งหมดที่ Admin อาจต้องการดู
	result := r.db.
		Preload("Department").
		Preload("Role").
		Preload("Major").
		Preload("Prefix").
		Preload("ManagedUsers.Role").
		Preload("ManagedUsers.Prefix").
		Preload("AdvisorProfile").
		Preload("StudentProfile").
		Preload("StudentProfile.StudentAcademicRecords", func(db *gorm.DB) *gorm.DB {
			return db.Order("academic_year desc, semester desc")
		}).
		Where("sut_id = ?", sutID).
		First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *AdminProfileRepositoryImpl) FindAllMajors() ([]*entity.Major, error) {
	var majors []*entity.Major
	// สมมติว่า Major table เป็นตาราง Master Data
	result := r.db.Find(&majors)
	if result.Error != nil {
		return nil, result.Error
	}
	return majors, nil
}

func (r *AdminProfileRepositoryImpl) UpdateUserStatus(sutID string, newStatus string) error {
	var newActiveBool bool
	switch newStatus {
	case "active":
		newActiveBool = true
	case "inactive":
		newActiveBool = false
	}
	result := r.db.Model(&entity.User{}).
		Where("sut_id = ?", sutID).
		Update("active", newActiveBool)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *AdminProfileRepositoryImpl) FindUserBySutID(sutID string) (*entity.User, error) {
	var user entity.User

	if err := r.db.Where("sut_id = ?", sutID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}
