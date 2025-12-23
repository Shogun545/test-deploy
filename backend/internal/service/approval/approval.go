package service

import (
	"backend/internal/app/entity"
	"backend/internal/app/repository"
	"errors"
	"time"
)

const (
	StatusPending    uint = 1
	StatusApproved   uint = 2
	StatusReschedule uint = 3

	ActionApprove    uint = 1
	ActionReschedule uint = 2
)

type AppointmentService interface {
	ApproveAppointment(AppointmentID uint, ActorID uint, Role string, Description string) (*entity.Appointment, error)
	ProposeNewTime(AppointmentID uint, ActorID uint, Role string, Description string) (*entity.Appointment, error)

	GetByID(id uint) (*entity.Appointment, error)
	ListPendingByAdvisor(advisorID uint) ([]entity.Appointment, error)

	// ✅ เพิ่มใหม่
	ListDoneByAdvisor(advisorID uint) ([]entity.Appointment, error)
	ListAllByAdvisor(advisorID uint) ([]entity.Appointment, error)
}

type appointmentService struct {
	repo repository.AppointmentRepository
}

func NewAppointmentService(repo repository.AppointmentRepository) AppointmentService {
	return &appointmentService{repo: repo}
}

func (s *appointmentService) ApproveAppointment(
	AppointmentID uint,
	ActorID uint,
	Role string,
	Description string,
) (*entity.Appointment, error) {

	appt, err := s.repo.GetByID(AppointmentID)
	if err != nil {
		return nil, err
	}

	if Role != "advisor" && Role != "ADVISOR" {
		return nil, errors.New("only advisor can approve appointment")
	}

	if appt.AdvisorUserID != ActorID {
		return nil, errors.New("you are not the advisor of this appointment")
	}

	if appt.AppointmentStatusID != StatusPending {
		return nil, errors.New("appointment is not in pending status")
	}

	oldStatus := appt.AppointmentStatusID

	updateFields := map[string]interface{}{
		"appointment_status_id": StatusApproved,
	}
	if Description != "" {
		updateFields["description"] = Description
	}

	if err := s.repo.UpdateFields(appt.ID, updateFields); err != nil {
		return nil, err
	}

	if err := s.repo.UpsertAppointmentState(appt.ID, StatusApproved); err != nil {
		return nil, err
	}

	history := &entity.StatusHistory{
		AppointmentID:   appt.ID,
		ChangedByUserID: ActorID,
		FromStatusID:    oldStatus,
		ToStatusID:      StatusApproved,
		ActionID:        ActionApprove,
		Reason:          Description,
		ChangedAt:       time.Now(),
	}
	if err := s.repo.CreateStatusHistory(history); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateFields(appt.ID, map[string]interface{}{
		"status_history_id": history.ID,
	}); err != nil {
		return nil, err
	}

	updated, _ := s.repo.GetByID(appt.ID)
	return updated, nil
}

func (s *appointmentService) ProposeNewTime(
	AppointmentID uint,
	ActorID uint,
	Role string,
	Description string,
) (*entity.Appointment, error) {

	appt, err := s.repo.GetByID(AppointmentID)
	if err != nil {
		return nil, err
	}

	if Role != "advisor" && Role != "ADVISOR" {
		return nil, errors.New("only teacher/advisor can propose new time")
	}

	if appt.AdvisorUserID != ActorID {
		return nil, errors.New("you are not the advisor of this appointment")
	}

	if appt.AppointmentStatusID != StatusPending {
		return nil, errors.New("appointment is not in pending status")
	}

	oldStatus := appt.AppointmentStatusID

	updateFields := map[string]interface{}{
		"appointment_status_id": StatusReschedule,
	}
	if Description != "" {
		updateFields["description"] = Description
	}

	if err := s.repo.UpdateFields(appt.ID, updateFields); err != nil {
		return nil, err
	}

	if err := s.repo.UpsertAppointmentState(appt.ID, StatusReschedule); err != nil {
		return nil, err
	}

	history := &entity.StatusHistory{
		AppointmentID:   appt.ID,
		ChangedByUserID: ActorID,
		FromStatusID:    oldStatus,
		ToStatusID:      StatusReschedule,
		ActionID:        ActionReschedule,
		Reason:          Description,
		ChangedAt:       time.Now(),
	}
	if err := s.repo.CreateStatusHistory(history); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateFields(appt.ID, map[string]interface{}{
		"status_history_id": history.ID,
	}); err != nil {
		return nil, err
	}

	updated, _ := s.repo.GetByID(appt.ID)
	return updated, nil
}

func (s *appointmentService) GetByID(id uint) (*entity.Appointment, error) {
	return s.repo.GetByID(id)
}

func (s *appointmentService) ListPendingByAdvisor(advisorID uint) ([]entity.Appointment, error) {
	return s.repo.ListPendingByAdvisor(advisorID)
}

// ✅ เพิ่มใหม่
func (s *appointmentService) ListDoneByAdvisor(advisorID uint) ([]entity.Appointment, error) {
	return s.repo.ListDoneByAdvisor(advisorID)
}

func (s *appointmentService) ListAllByAdvisor(advisorID uint) ([]entity.Appointment, error) {
	return s.repo.ListAllByAdvisor(advisorID)
}
