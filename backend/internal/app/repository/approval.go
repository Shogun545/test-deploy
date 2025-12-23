package repository

import (
	"backend/internal/app/entity"
	"errors"

	"gorm.io/gorm"
)

type AppointmentRepository interface {
	GetByID(id uint) (*entity.Appointment, error)
	Update(appt *entity.Appointment) error
	UpdateFields(id uint, fields map[string]interface{}) error

	UpsertAppointmentState(appointmentID uint, statusID uint) error
	CreateStatusHistory(h *entity.StatusHistory) error

	ListPendingByAdvisor(advisorID uint) ([]entity.Appointment, error)

	// ✅ เพิ่มใหม่
	ListDoneByAdvisor(advisorID uint) ([]entity.Appointment, error)
	ListAllByAdvisor(advisorID uint) ([]entity.Appointment, error)
}

type appointmentRepository struct {
	db *gorm.DB
}

func NewAppointmentRepository(db *gorm.DB) AppointmentRepository {
	return &appointmentRepository{db: db}
}

func (r *appointmentRepository) GetByID(id uint) (*entity.Appointment, error) {
	var appt entity.Appointment
	if err := r.db.
		Preload("AdvisorUser").
		Preload("StudentUser").
		Preload("Topic").
		Preload("Category").
		Preload("AppointmentStatus").
		First(&appt, id).Error; err != nil {
		return nil, err
	}
	return &appt, nil
}

func (r *appointmentRepository) Update(a *entity.Appointment) error {
	return r.db.Save(a).Error
}

func (r *appointmentRepository) UpdateFields(id uint, fields map[string]interface{}) error {
	return r.db.Model(&entity.Appointment{}).
		Where("id = ?", id).
		Updates(fields).Error
}

func (r *appointmentRepository) UpsertAppointmentState(appointmentID uint, statusID uint) error {
	var state entity.AppointmentState

	err := r.db.Where("appointment_id = ?", appointmentID).First(&state).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			state = entity.AppointmentState{
				AppointmentID:   appointmentID,
				CurrentStatusID: statusID,
			}
			return r.db.Create(&state).Error
		}
		return err
	}

	state.CurrentStatusID = statusID
	return r.db.Save(&state).Error
}

func (r *appointmentRepository) CreateStatusHistory(h *entity.StatusHistory) error {
	return r.db.Create(h).Error
}

func (r *appointmentRepository) ListPendingByAdvisor(advisorID uint) ([]entity.Appointment, error) {
	var appts []entity.Appointment
	err := r.db.
		Preload("StudentUser").
		Preload("Topic").
		Preload("AppointmentStatus").
		Where("advisor_user_id = ? AND appointment_status_id = ?", advisorID, entity.StatusPendingID).
		Order("created_at DESC").
		Find(&appts).Error
	return appts, err
}

// ✅ พิจารณาแล้ว = Approved + Reschedule
func (r *appointmentRepository) ListDoneByAdvisor(advisorID uint) ([]entity.Appointment, error) {
	var appts []entity.Appointment
	err := r.db.
		Preload("StudentUser").
		Preload("Topic").
		Preload("AppointmentStatus").
		Where("advisor_user_id = ? AND appointment_status_id IN ?", advisorID, []uint{entity.StatusApprovedID, entity.StatusRescheduleID}).
		Order("created_at DESC").
		Find(&appts).Error
	return appts, err
}

// ✅ ทั้งหมด = ไม่กรอง status
func (r *appointmentRepository) ListAllByAdvisor(advisorID uint) ([]entity.Appointment, error) {
	var appts []entity.Appointment
	err := r.db.
		Preload("StudentUser").
		Preload("Topic").
		Preload("AppointmentStatus").
		Where("advisor_user_id = ?", advisorID).
		Order("created_at DESC").
		Find(&appts).Error
	return appts, err
}
