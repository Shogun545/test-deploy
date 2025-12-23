package advisorprofile

import (
	"backend/internal/app/dto"
	//"backend/internal/app/entity"
	"backend/internal/app/repository"
	"errors"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("record not found")
var ErrUpdateFailed = errors.New("failed to update profile")

type AdvisorProfileService struct {
	Repo repository.AdvisorProfileRepository
	db   *gorm.DB
}

func NewAdvisorProfileService(repo repository.AdvisorProfileRepository, db *gorm.DB) *AdvisorProfileService {
	return &AdvisorProfileService{
		Repo: repo,
		db:   db,
	}
}

// สร้าง method GetMyProfile
func (s *AdvisorProfileService) GetMyProfile(sutID string) (*dto.AdvisorProfileResponse, error) {
	profile, err := s.Repo.FindBySutID(sutID)
	if err != nil {
		return nil, ErrNotFound
	}
	resp := dto.NewAdvisorProfileResponse(profile.User, profile)
	return &resp, nil
}

// method UpdateMyProfile
func (s *AdvisorProfileService) UpdateMyProfile(
	sutID string,
	req *dto.UpdateAdvisorProfileRequest,
) (*dto.AdvisorProfileResponse, error) {

	// 1. ดึง Advisor พร้อมโหลดข้อมูล User และ Prefix มาด้วย
	advisor, err := s.Repo.FindBySutID(sutID) // ตรวจสอบให้แน่ใจว่า Repo นี้ Preload("User") แล้ว
	if err != nil {
		return nil, ErrNotFound
	}

	user := advisor.User // ดึง User entity ออกมา

	// 2. จัดการคำนำหน้า (Prefix)
	if req.Prefix != "" {
		prefix, err := s.Repo.FindPrefixByName(req.Prefix)
		if err == nil {
			user.PrefixID = prefix.ID
			user.Prefix = prefix // อัปเดตเพื่อให้ Response มีข้อมูลใหม่ทันที
		}
	}

	// 3. อัปเดต ชื่อ-นามสกุล และข้อมูลอื่นๆ ของ User
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.ProfileImage != "" {
		user.ProfileImage = req.ProfileImage
	}

	// 4. อัปเดตข้อมูลของ Advisor Profile
	if req.OfficeRoom != "" {
		advisor.OfficeRoom = req.OfficeRoom
	}
	if req.Specialties != "" {
		advisor.Specialties = req.Specialties
	}
	if req.IsActive != nil {
		advisor.IsActive = *req.IsActive
	}

	// 5. บันทึกข้อมูลแบบ Transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// เซฟ User ก่อน (สำหรับชื่อและ Prefix)
		if err := tx.Model(user).Select(
			"PrefixID",
			"FirstName",
			"LastName",
			"Phone",
			"Email",
			"ProfileImage", 
		).Updates(user).Error; err != nil {
			return err
		}

		// เซฟ Advisor Profile (สำหรับห้องทำงานและ isActive)
		if err := tx.Save(advisor).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, ErrUpdateFailed
	}

	// 6. สร้าง Response ส่งกลับ
	resp := dto.NewAdvisorProfileResponse(user, advisor)
	return &resp, nil
}

func (s *AdvisorProfileService) GetMyStudents(sutID string) (*dto.AdvisorStudentsResponse, error) {
	advisor, err := s.Repo.FindAdvisorWithStudentsBySutID(sutID)
	if err != nil {
		return nil, ErrNotFound
	}

	resp := dto.NewAdvisorStudentsResponse(advisor)
	return &resp, nil
}

func (s *AdvisorProfileService) GetStudentBySutID(
	advisorSutID string,
	studentSutID string,
) (*dto.AdvisorStudentDetailResponse, error) {

	advisor, err := s.Repo.FindAdvisorWithStudentsBySutID(advisorSutID)
	if err != nil {
		return nil, ErrNotFound
	}

	for _, st := range advisor.AdvisorProfile.Students {
		if st.User.SutId == studentSutID {

			major := ""
			if st.User.Major != nil {
				major = st.User.Major.Major
			}

			year := 0
			gpa := 0.0
			if st.User.StudentProfile != nil {
				year = st.User.StudentProfile.YearOfStudy
				if len(st.User.StudentProfile.StudentAcademicRecords) > 0 {
					gpa = float64(
						st.User.StudentProfile.StudentAcademicRecords[0].CumulativeGPA,
					)
				}
			}

			return &dto.AdvisorStudentDetailResponse{
				SutId:       st.User.SutId,
				FirstName:   st.User.FirstName,
				LastName:    st.User.LastName,
				YearOfStudy: year,
				Email:       st.User.Email,
				Phone:       st.User.Phone,
				MajorName:   major,
				GPALatest:   gpa,
				Birthday:    st.User.Birthday,
			}, nil
		}
	}

	return nil, ErrNotFound
}
