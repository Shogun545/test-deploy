package adminprofile

import (
	"backend/internal/app/dto"
	"backend/internal/app/entity"
	"backend/internal/app/repository"
	"errors"
	"gorm.io/gorm"
)

// กำหนด Error

var ErrNotFound = errors.New("record not found")
var ErrUpdateFailed = errors.New("failed to update record")
var ErrInvalidStatus = errors.New("invalid user status value")

// AdminProfileService ใช้งาน repository ที่สามารถจัดการ User/Admin ได้
type AdminProfileService struct {
	Repo repository.AdminProfileRepository // Repository สำหรับ Admin
	db   *gorm.DB
}

func NewAdminProfileService(repo repository.AdminProfileRepository, db *gorm.DB) *AdminProfileService {
	return &AdminProfileService{
		Repo: repo,
		db:   db,
	}
}

// GetMyProfile ดึงข้อมูลโปรไฟล์ของ Admin ที่กำลังล็อกอินอยู่
func (s *AdminProfileService) GetMyProfile(adminSutID string) (*dto.AdminProfileResponse, error) {
	// สมมติว่า FindBySutID ใน AdminProfileRepository ดึงข้อมูล User (Admin) พร้อม Department/Role
	profile, err := s.Repo.FindBySutID(adminSutID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// ต้องสร้าง DTO Response สำหรับ Admin โดยเฉพาะ
	resp := dto.NewAdminProfileResponse(profile)
	return &resp, nil
}

// UpdateMyProfile อัปเดตข้อมูลส่วนตัวของ Admin (Controller ที่ถูกเพิ่มในภายหลังจะเรียกใช้)
func (s *AdminProfileService) UpdateMyProfile(
	adminSutID string,
	req *dto.UpdateAdminProfileRequest,
) (*dto.AdminProfileResponse, error) {

	// ดึง User (Admin) ปัจจุบัน
	var admin entity.User
	if err := s.db.Where("sut_id = ?", adminSutID).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if req.Prefix != "" {
		var prefix entity.Prefix
		if err := s.db.Where("prefix = ?", req.Prefix).First(&prefix).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("invalid prefix")
			}
			return nil, err
		}
		admin.PrefixID = prefix.ID
	}
	if req.FirstName != "" {
		admin.FirstName = req.FirstName
	}
	if req.LastName != "" {
		admin.LastName = req.LastName
	}

	// อัปเดตฟิลด์ที่ได้รับจาก Request
	if req.Phone != "" {
		admin.Phone = req.Phone
	}
	if req.Notes != "" {
		admin.Notes = req.Notes
	}
	if req.Active != nil {
		admin.Active = *req.Active // สมมติว่า Active เป็น bool
	}
	if req.Email != "" {
		admin.Email = req.Email
	}

	// บันทึกการเปลี่ยนแปลง
	if err := s.db.Save(&admin).Error; err != nil {
		return nil, ErrUpdateFailed
	}

	// ดึงข้อมูลที่อัปเดตแล้วและสร้าง Response (อาจต้อง Preload Department อีกครั้ง)
	// เพื่อให้ได้ข้อมูลสมบูรณ์สำหรับ Response
	updatedAdmin, err := s.Repo.FindBySutID(adminSutID)
	if err != nil {
		return nil, err
	}

	resp := dto.NewAdminProfileResponse(updatedAdmin)
	return &resp, nil
}

// GetAllUsers ดึงข้อมูลผู้ใช้งานทั้งหมดในระบบ (เพื่อแสดงในหน้า Managed Users)
// รับ Struct UserListFilter
func (s *AdminProfileService) GetAllUsers(filters dto.UserListFilter) (*dto.AdminManagedUsersResponse, error) {

	//  ต้องไปแก้ไข Interface และ Implementation ของ FindAllUsers ใน Repository ด้วย
	users, err := s.Repo.FindAllUsers(filters)
	if err != nil {
		return nil, err
	}

	resp := dto.NewAdminManagedUsersResponse(users)
	return &resp, nil
}

// GetUserBySutID ดึงรายละเอียดผู้ใช้งานคนใดคนหนึ่ง (ทั้งนักศึกษาและอาจารย์)
func (s *AdminProfileService) GetUserBySutID(targetSutID string) (*dto.AdminUserDetailResponse, error) {
	// Repo.FindUserDetailBySutID ดึง User พร้อมข้อมูลที่เกี่ยวข้องทั้งหมด (Profile, Academic, Role)
	user, err := s.Repo.FindUserDetailBySutID(targetSutID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// ต้องสร้าง DTO Response สำหรับแสดงรายละเอียด User ที่ Admin ดู
	resp := dto.NewAdminUserDetailResponse(user)
	return &resp, nil
}

func (s *AdminProfileService) GetMasterMajors() ([]dto.MajorEntry, error) {
	majors, err := s.Repo.FindAllMajors()
	if err != nil {
		return nil, err
	}

	resp := make([]dto.MajorEntry, len(majors))
	for i, m := range majors {
		resp[i] = dto.NewMajorEntry(m)
	}
	return resp, nil
}

func (s *AdminProfileService) UpdateUserStatus(sutID string, newStatus string) error {
	// ...
	// 2. เรียก Repository เพื่ออัปเดตข้อมูลใน DB
	err := s.Repo.UpdateUserStatus(sutID, newStatus)
	if err != nil {
		// จัดการ Error จาก DB (เช่น ErrNotFound)
		// เปลี่ยนจาก errors.Is(err, repo.ErrNotFound) เป็น:
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func (s *AdminProfileService) GetUserCreatedDate(sutID string) (*dto.UserCreatedDateResponse, error) {
	user, err := s.Repo.FindUserBySutID(sutID)
	if err != nil {
		return nil, ErrNotFound
	}

	return &dto.UserCreatedDateResponse{
		SutID:     user.SutId,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *AdminProfileService) UpdateUser(
	sutID string,
	req *dto.UpdateManagedUserRequest,
) error {

	user, err := s.Repo.FindUserBySutID(sutID)
	if err != nil {
		return ErrNotFound
	}

	if req.Phone != nil {
		user.Phone = *req.Phone
	}

	if err := s.Repo.UpdateUser(user); err != nil {
		return ErrUpdateFailed
	}

	return nil
}
