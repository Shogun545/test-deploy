package service

import (
	"backend/internal/app/dto"
	"backend/internal/app/entity"
	"backend/internal/app/repository"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInternal     = errors.New("internal error")
	ErrUnauthorized = errors.New("unauthorized")
)

type ProfileService struct {
	Repo *repository.ProfileRepository
	DB   *gorm.DB
}

func NewProfileService(repo *repository.ProfileRepository) *ProfileService {
	return &ProfileService{
		Repo: repo,
		DB:   repo.DB,
	}
}

// คืน dto.NewStudentProfileResponse หรือ error (ErrNotFound / ErrInternal)
func (s *ProfileService) GetMyProfile(sutID string) (*dto.StudentProfileResponse, error) {
	if sutID == "" {
		return nil, ErrUnauthorized
	}

	user, err := s.Repo.FindUserWithProfileBySutID(sutID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}

	sp := user.StudentProfile
	if sp == nil {
		return nil, ErrNotFound
	}

	var record entity.StudentAcademicRecord
	recordPtr, err := s.Repo.FindLatestAcademicRecord(sp.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInternal
	}
	if recordPtr != nil {
		record = *recordPtr
	}

	termLabel := ""
	if record.AcademicYear != "" {
		termLabel = fmt.Sprintf("ภาคการศึกษา %d/%s", record.Semester, record.AcademicYear)
	}

	advisorName := ""
	if sp.AdvisorProfile != nil && sp.AdvisorProfile.User != nil {
		advisorName = sp.AdvisorProfile.User.FirstName + " " + sp.AdvisorProfile.User.LastName
	}

	// โหลด user ใหม่ เพื่อให้ Prefix ถูกต้อง
	updatedUser, err := s.Repo.FindUserWithProfileBySutID(sutID)
	if err != nil {
		return nil, ErrInternal
	}

	resp := dto.NewStudentProfileResponse(
		updatedUser,
		updatedUser.StudentProfile,
		&record,
		advisorName,
		termLabel,
	)

	return &resp, nil

}

// Update: รับ sutID และ DTO request → อัปเดต fields ที่อนุญาต แล้วคืน StudentProfileResponse
func (s *ProfileService) UpdateMyProfile(
	sutID string,
	req *dto.UpdateStudentProfileRequest,
) (*dto.StudentProfileResponse, error) {

	if sutID == "" {
		return nil, ErrUnauthorized
	}

	// 1. โหลดข้อมูลเดิม
	user, err := s.Repo.FindUserWithProfileBySutID(sutID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}

	sp := user.StudentProfile
	if sp == nil {
		return nil, ErrNotFound
	}

	// 2. อัปเดตข้อมูลส่วนตัว
	if req.Prefix != "" {
		prefix, err := s.Repo.FindPrefixByName(req.Prefix)
		if err == nil {
			user.PrefixID = prefix.ID
			user.Prefix = prefix
		}
	}
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Birthday != "" {
		user.Birthday = req.Birthday
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.ProfileImage != "" {
		user.ProfileImage = req.ProfileImage
	}

	if req.YearOfStudy != nil {
		sp.YearOfStudy = *req.YearOfStudy
	}

	// 3. จัดการเรื่องเกรดและสถานะ (Academic Record)
	record, _ := s.Repo.FindLatestAcademicRecord(sp.ID)
	if record != nil {
		if req.AcademicYear != "" {
			record.AcademicYear = req.AcademicYear
		} // อัปเดตปี
		if req.Semester != nil {
			record.Semester = *req.Semester
		} // อัปเดตเทอม
		if req.CumulativeGPA != nil {
			record.CumulativeGPA = *req.CumulativeGPA
		}
		if req.TermGPA != nil {
			record.TermGPA = *req.TermGPA
		}

		// Logic คำนวณสถานะ (เชื่อมโยงตามเกรดสะสม)
		if record.CumulativeGPA < 1.50 {
			record.AcademicStatus = "พ้นสภาพ"
		} else if record.CumulativeGPA < 2.00 {
			record.AcademicStatus = "รอพินิจ (Probation)"
		} else {
			record.AcademicStatus = "ปกติ"
		}
	}

	// 4. บันทึกลง Database ภายใน Transaction เดียว
	err = s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(user).Select("PrefixID", "FirstName", "LastName", "Phone", "Birthday", "Email", "ProfileImage").Updates(user).Error; err != nil {
			return err
		}
		if err := tx.Model(sp).Select("YearOfStudy").Updates(sp).Error; err != nil {
			return err
		}
		if record != nil {
			if err := tx.Save(record).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, ErrInternal
	}

	// 5. โหลดข้อมูลใหม่ทั้งหมดเพื่อส่งกลับ
	updatedUser, _ := s.Repo.FindUserWithProfileBySutID(sutID)
	latestRec, _ := s.Repo.FindLatestAcademicRecord(sp.ID)

	termLabel := ""
	advisorName := ""
	if latestRec != nil {
		termLabel = fmt.Sprintf("ภาคการศึกษา %d/%s", latestRec.Semester, latestRec.AcademicYear)
	}
	if updatedUser.StudentProfile.AdvisorProfile != nil && updatedUser.StudentProfile.AdvisorProfile.User != nil {
		advisorName = updatedUser.StudentProfile.AdvisorProfile.User.FirstName + " " + updatedUser.StudentProfile.AdvisorProfile.User.LastName
	}

	resp := dto.NewStudentProfileResponse(updatedUser, updatedUser.StudentProfile, latestRec, advisorName, termLabel)
	return &resp, nil
}
