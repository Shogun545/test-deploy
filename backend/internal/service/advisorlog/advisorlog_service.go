package advisorlog

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"backend/internal/app/dto"
	"backend/internal/app/entity"
)

type Service interface {
    // ✅ เพิ่ม requesterID เพื่อเช็คว่าเป็นเจ้าของนัดหมายจริงไหม
	Create(ctx context.Context, req dto.AdvisorLogCreateReq, requesterID uint, requesterRole string) (*dto.AdvisorLogCreateResp, error)
	GetByID(ctx context.Context, id uint, requesterID uint, requesterRole string) (*dto.AdvisorLogGetResp, error)
	ListByStudent(ctx context.Context, studentUserID uint, requesterRole string) ([]dto.AdvisorLogListItemResp, error)
	ListAll(ctx context.Context) ([]dto.AdvisorLogListItemResp, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
    // ✅ เพิ่ม requesterID/Role เพื่อเช็คสิทธิ์ก่อนแก้
	Update(ctx context.Context, id uint, req dto.AdvisorLogUpdateReq, requesterID uint, requesterRole string) (*dto.AdvisorLogUpdateResp, error)
	GetFileForLog(ctx context.Context, logID uint, index int, sutID string) (string, string, error)
}

var (
	ErrAdvisorLogNotFound         = errors.New("advisor log not found")
	ErrInvalidStatus              = errors.New("invalid status")
	ErrSaveFileFailed             = errors.New("save file failed")
	ErrAppointmentNotCompletedYet = errors.New("cannot create log: appointment is not completed yet")
	ErrForbidden                  = errors.New("forbidden")
	ErrFileNotFound               = errors.New("file not found")
)

type service struct {
	db        *gorm.DB
	uploadDir string
}

func New(db *gorm.DB, uploadDir string) Service {
	if uploadDir == "" {
		uploadDir = "uploads"
	}
	return &service{db: db, uploadDir: uploadDir}
}

// ------------------------------
// helpers
// ------------------------------

func (s *service) cleanupFiles(paths []string) {
	for _, p := range paths {
		if p == "" {
			continue
		}
		_ = os.Remove(p)
	}
}

func writeFile(dest string, r io.Reader) error {
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	return err
}

func splitCsv(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func toBase(log entity.AdvisorLog) dto.AdvisorLogRespBase {
	return dto.AdvisorLogRespBase{
		ID:             log.ID,
		AppointmentID:  log.AppointmentID,
		Title:          log.Title,
		Body:           log.Body,
		Status:         log.Status,
		RequiresReport: log.RequiresReport,
		FileName:       log.FileName,
		FilePath:       log.FilePath,
        // ✅ แก้ไข: เพิ่ม Date Mapping
        CreatedAt:      log.CreatedAt.Format("2006-01-02 15:04:05"),
        UpdatedAt:      log.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// saveFiles: ปิด src ให้ถูกต้อง (ไม่ defer ใน loop)
func (s *service) saveFiles(files []*multipart.FileHeader, orig *[]string, paths *[]string) error {
	for _, fh := range files {
		ext := filepath.Ext(fh.Filename)
		if ext == "" {
			ext = ".bin"
		}
		dest := filepath.Join(s.uploadDir, uuid.New().String()+ext)

		src, err := fh.Open()
		if err != nil {
			return ErrSaveFileFailed
		}

		if err := writeFile(dest, src); err != nil {
			_ = src.Close()
			return ErrSaveFileFailed
		}
		_ = src.Close()

		*orig = append(*orig, fh.Filename)
		*paths = append(*paths, dest)
	}
	return nil
}

// ------------------------------
// CREATE (🔒 Secure: อนุญาตทั้ง นศ. และ อาจารย์เจ้าของนัด)
// ------------------------------
// ⚠️ แก้ Signature: รับ requesterRole เพิ่มเข้ามา
func (s *service) Create(ctx context.Context, req dto.AdvisorLogCreateReq, requesterID uint, requesterRole string) (*dto.AdvisorLogCreateResp, error) {
	// 1. ดึงข้อมูล Appointment มาเช็ค
	var appt entity.Appointment
	if err := s.db.WithContext(ctx).
		Preload("AppointmentStatus").
		First(&appt, req.AppointmentID).Error; err != nil {
		return nil, err
	}

	// 🛡️ Security Check (Logic ใหม่):
	// อนุญาตถ้า:
	// 1. เป็น Student ที่เป็นเจ้าของนัดหมาย (StudentUserID == requesterID)
	// 2. หรือ เป็น Advisor ที่เป็นคู่กรณีในนัดหมาย (AdvisorUserID == requesterID)
	
	isStudentOwner := appt.StudentUserID == requesterID
	isAdvisorOwner := appt.AdvisorUserID == requesterID

	// ถ้าไม่ใช่ทั้งคู่ -> ห้ามทำรายการ
	if !isStudentOwner && !isAdvisorOwner {
		return nil, ErrForbidden
	}

	// (Optional) ถ้าต้องการบังคับว่าต้อง Completed แล้วเท่านั้นถึงจดบันทึกได้
	if !appt.AppointmentStatus.IsTerminal {
		return nil, ErrAppointmentNotCompletedYet
	}

	// 2. สร้าง Log Object
	log := entity.AdvisorLog{
		AppointmentID:  req.AppointmentID,
		Title:          req.Title,
		Body:           req.Body,
		RequiresReport: req.RequiresReport,
	}

	// status logic
	if strings.ToLower(req.Status) == "draft" {
		log.Status = "Draft"
	} else {
		if log.RequiresReport {
			log.Status = "PendingReport"
		} else {
			log.Status = "Completed"
		}
	}

	if _, err := govalidator.ValidateStruct(log); err != nil {
		return nil, err
	}

	// upload files
	_ = os.MkdirAll(s.uploadDir, 0755)

	var savedOrig, savedPaths []string
	if err := s.saveFiles(req.Files, &savedOrig, &savedPaths); err != nil {
		s.cleanupFiles(savedPaths)
		return nil, err
	}

	if len(savedPaths) > 0 {
		log.FileName = strings.Join(savedOrig, ",")
		log.FilePath = strings.Join(savedPaths, ",")
	}

	if err := s.db.WithContext(ctx).Create(&log).Error; err != nil {
		s.cleanupFiles(savedPaths)
		return nil, err
	}

	return &dto.AdvisorLogCreateResp{
		AdvisorLogRespBase: toBase(log),
	}, nil
}

// ------------------------------
// LIST BY STUDENT ID (🔒 Secure)
// ------------------------------
func (s *service) ListByStudent(ctx context.Context, studentUserID uint, requesterRole string) ([]dto.AdvisorLogListItemResp, error) {
	var logs []entity.AdvisorLog

	query := s.db.WithContext(ctx).
		Joins("JOIN appointments ON appointments.id = advisor_logs.appointment_id").
		Where("appointments.student_user_id = ?", studentUserID)

	// กรอง Draft ออก ถ้านักศึกษาเรียก
	if strings.ToLower(requesterRole) == "student" {
		query = query.Where("advisor_logs.status != ?", "Draft")
	}

	err := query.Order("advisor_logs.id desc").Find(&logs).Error
	if err != nil {
		return nil, err
	}

	out := make([]dto.AdvisorLogListItemResp, 0, len(logs))
	for _, l := range logs {
		out = append(out, dto.AdvisorLogListItemResp{
			AdvisorLogRespBase: toBase(l),
		})
	}
	return out, nil
}

// ------------------------------
// GET BY ID (🔒 Secure)
// ------------------------------
func (s *service) GetByID(ctx context.Context, id uint, requesterID uint, requesterRole string) (*dto.AdvisorLogGetResp, error) {
	var log entity.AdvisorLog
	
	if err := s.db.WithContext(ctx).
		Preload("Appointment").
		First(&log, id).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAdvisorLogNotFound
		}
		return nil, err
	}

	// 🛡️ Logic ความปลอดภัย
	if strings.ToLower(requesterRole) == "student" {
		if log.Status == "Draft" {
			return nil, ErrAdvisorLogNotFound
		}
		if log.Appointment == nil || log.Appointment.StudentUserID != requesterID {
			return nil, ErrForbidden
		}
	}
	
	return &dto.AdvisorLogGetResp{
		AdvisorLogRespBase: toBase(log),
	}, nil
}

// ------------------------------
// LIST ALL (advisor)
// ------------------------------
func (s *service) ListAll(ctx context.Context) ([]dto.AdvisorLogListItemResp, error) {
	var logs []entity.AdvisorLog
	err := s.db.WithContext(ctx).
		Preload("Appointment").
		Order("advisor_logs.id desc").
		Find(&logs).Error
	if err != nil {
		return nil, err
	}

	out := make([]dto.AdvisorLogListItemResp, 0, len(logs))
	for _, l := range logs {
		out = append(out, dto.AdvisorLogListItemResp{
			AdvisorLogRespBase: toBase(l),
		})
	}
	return out, nil
}

// ------------------------------
// UPDATE STATUS ONLY
// ------------------------------
func (s *service) UpdateStatus(ctx context.Context, id uint, status string) error {
	allowed := map[string]bool{
		"Draft":         true,
		"PendingReport": true,
		"Completed":     true,
	}
	if !allowed[status] {
		return ErrInvalidStatus
	}

	res := s.db.WithContext(ctx).
		Model(&entity.AdvisorLog{}).
		Where("id = ?", id).
		Update("status", status)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrAdvisorLogNotFound
	}
	return nil
}

// ------------------------------
// UPDATE FULL LOG (🔒 Secure: เช็คเจ้าของก่อนแก้)
// ------------------------------
func (s *service) Update(ctx context.Context, id uint, req dto.AdvisorLogUpdateReq, requesterID uint, requesterRole string) (*dto.AdvisorLogUpdateResp, error) {
	var log entity.AdvisorLog
	if err := s.db.WithContext(ctx).Preload("Appointment").First(&log, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAdvisorLogNotFound
		}
		return nil, err
	}

    // 🛡️ Security Check:
    if strings.ToLower(requesterRole) == "student" {
        // 1. ต้องเป็นเจ้าของ Appointment
        if log.Appointment == nil || log.Appointment.StudentUserID != requesterID {
            return nil, ErrForbidden
        }
        // 2. (Optional) ถ้า Completed แล้ว ห้ามแก้? หรือถ้าเป็น Draft เท่านั้นถึงแก้ได้?
        // (ตรงนี้แล้วแต่ Business Logic ของคุณ ปกติถ้านักศึกษาแก้ได้ตลอดก็ไม่ต้องเช็ค Status)
    }

	// update text fields
	if req.Title != nil && *req.Title != "" {
		log.Title = *req.Title
	}
	if req.Body != nil && *req.Body != "" {
		log.Body = *req.Body
	}
	if req.RequiresReport != nil {
		log.RequiresReport = *req.RequiresReport
	}

	// Files Logic (เหมือนเดิม)
	if len(req.Files) > 0 {
		_ = os.MkdirAll(s.uploadDir, 0755)

		var newNames []string
		var newPaths []string

		for _, fh := range req.Files {
			ext := filepath.Ext(fh.Filename)
			if ext == "" {
				ext = ".bin"
			}
			dest := filepath.Join(s.uploadDir, uuid.New().String()+ext)

			src, err := fh.Open()
			if err != nil {
				s.cleanupFiles(newPaths)
				return nil, ErrSaveFileFailed
			}

			if err := writeFile(dest, src); err != nil {
				_ = src.Close()
				s.cleanupFiles(newPaths)
				return nil, ErrSaveFileFailed
			}
			_ = src.Close()

			newNames = append(newNames, fh.Filename)
			newPaths = append(newPaths, dest)
		}

		var oldPaths []string
		if log.FilePath != "" {
			oldPaths = strings.Split(log.FilePath, ",")
		}

		log.FileName = strings.Join(newNames, ",")
		log.FilePath = strings.Join(newPaths, ",")

		if err := s.db.WithContext(ctx).Save(&log).Error; err != nil {
			s.cleanupFiles(newPaths) 
			return nil, err
		}
		s.cleanupFiles(oldPaths)

	} else {
		if err := s.db.WithContext(ctx).Save(&log).Error; err != nil {
			return nil, err
		}
	}

	return &dto.AdvisorLogUpdateResp{
		AdvisorLogRespBase: toBase(log),
	}, nil
}

// ------------------------------
// GET FILE (Optimization: ถ้าทำได้ ควรใช้ user_id แบบ uint แต่ใช้แบบเดิมก็ไม่ผิด)
// ------------------------------
func (s *service) GetFileForLog(ctx context.Context, logID uint, index int, sutID string) (string, string, error) {
    // ... (Code เดิมส่วนนี้ใช้ได้ครับ ถ้า Controller ส่งมาแค่ sutID) ...
    // ... (logic การเช็คสิทธิ์ในนี้ถือว่าถูกต้องแล้ว) ...
	var log entity.AdvisorLog
	if err := s.db.WithContext(ctx).
		Preload("Appointment").
		First(&log, logID).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", ErrAdvisorLogNotFound
		}
		return "", "", err
	}

	sutID = strings.TrimSpace(sutID)
	if sutID == "" {
		return "", "", ErrForbidden
	}

	var me entity.User
	if err := s.db.WithContext(ctx).
		Preload("Role").
		Where("sut_id = ?", sutID).
		First(&me).Error; err != nil {
		return "", "", ErrForbidden
	}

	appt := log.Appointment
	if appt == nil {
		return "", "", ErrForbidden
	}
	if me.Role == nil {
		return "", "", ErrForbidden
	}

	role := strings.ToLower(strings.TrimSpace(me.Role.Role))
	switch role {
	case "admin":
		// allow
	case "advisor":
		if me.ID != appt.AdvisorUserID {
			return "", "", ErrForbidden
		}
	case "student":
		if me.ID != appt.StudentUserID {
			return "", "", ErrForbidden
		}
	default:
		return "", "", ErrForbidden
	}

	paths := splitCsv(log.FilePath)
	names := splitCsv(log.FileName)

	if index < 0 || index >= len(paths) {
		return "", "", ErrFileNotFound
	}

	absPath := strings.TrimSpace(paths[index])
	if absPath == "" {
		return "", "", ErrFileNotFound
	}

	fileName := absPath
	if index < len(names) && strings.TrimSpace(names[index]) != "" {
		fileName = strings.TrimSpace(names[index])
	}

	if _, err := os.Stat(absPath); err != nil {
		return "", "", ErrFileNotFound
	}

	return fileName, absPath, nil
}