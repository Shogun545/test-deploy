package test

import (
	"errors"
	"testing"
	"time"

	"backend/internal/app/entity"
	"backend/internal/app/repository"
	approval "backend/internal/service/approval"

	. "github.com/onsi/gomega"
)

// --------------------
// Fake Repository (Mock)
// --------------------
type fakeAppointmentRepo struct {
	appt *entity.Appointment

	// knobs for forcing errors
	getErr           error
	updateFieldsErr  error
	upsertStateErr   error
	createHistoryErr error

	// record calls
	lastUpdateID     uint
	lastUpdateFields map[string]interface{}
	updateCalls      []map[string]interface{} // ✅ เก็บทุกครั้งที่ UpdateFields ถูกเรียก

	lastUpsertID     uint
	lastUpsertStatus uint
	lastHistory      *entity.StatusHistory

	// simulate auto-increment history ID
	historySeq uint
}

var _ repository.AppointmentRepository = (*fakeAppointmentRepo)(nil)

func newFakeRepo(appt *entity.Appointment) *fakeAppointmentRepo {
	return &fakeAppointmentRepo{
		appt:       appt,
		historySeq: 100,
	}
}

func (f *fakeAppointmentRepo) GetByID(id uint) (*entity.Appointment, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	if f.appt == nil || f.appt.ID != id {
		return nil, errors.New("not found")
	}
	return f.appt, nil
}

func (f *fakeAppointmentRepo) Update(appt *entity.Appointment) error {
	// not used by service, but required by interface
	f.appt = appt
	return nil
}

func (f *fakeAppointmentRepo) UpdateFields(id uint, fields map[string]interface{}) error {
	if f.updateFieldsErr != nil {
		return f.updateFieldsErr
	}
	f.lastUpdateID = id

	// ✅ copy fields เพื่อกันโดนแก้ทับ และเก็บทุก call
	copied := map[string]interface{}{}
	for k, v := range fields {
		copied[k] = v
	}
	f.updateCalls = append(f.updateCalls, copied)
	f.lastUpdateFields = copied

	// simulate DB update
	if f.appt != nil && f.appt.ID == id {
		if v, ok := fields["appointment_status_id"]; ok {
			if s, ok2 := v.(uint); ok2 {
				f.appt.AppointmentStatusID = s
				switch s {
				case approval.StatusPending:
					f.appt.AppointmentStatus.StatusCode = "PENDING"
				case approval.StatusApproved:
					f.appt.AppointmentStatus.StatusCode = "APPROVED"
				case approval.StatusReschedule:
					f.appt.AppointmentStatus.StatusCode = "RESCHEDULE"
				}
			}
		}
		if v, ok := fields["description"]; ok {
			if s, ok2 := v.(string); ok2 {
				f.appt.Description = s
			}
		}
	}

	return nil
}

func (f *fakeAppointmentRepo) UpsertAppointmentState(appointmentID uint, statusID uint) error {
	if f.upsertStateErr != nil {
		return f.upsertStateErr
	}
	f.lastUpsertID = appointmentID
	f.lastUpsertStatus = statusID
	return nil
}

func (f *fakeAppointmentRepo) CreateStatusHistory(h *entity.StatusHistory) error {
	if f.createHistoryErr != nil {
		return f.createHistoryErr
	}
	f.historySeq++
	h.ID = f.historySeq
	f.lastHistory = h
	return nil
}

func (f *fakeAppointmentRepo) ListPendingByAdvisor(advisorID uint) ([]entity.Appointment, error) {
	return nil, nil
}
func (f *fakeAppointmentRepo) ListDoneByAdvisor(advisorID uint) ([]entity.Appointment, error) {
	return nil, nil
}
func (f *fakeAppointmentRepo) ListAllByAdvisor(advisorID uint) ([]entity.Appointment, error) {
	return nil, nil
}

// --------------------
// Tests: ApproveAppointment
// --------------------
func TestApproveAppointment(t *testing.T) {
	RegisterTestingT(t)

	t.Run("Case 1: Success - advisor approves pending appointment", func(t *testing.T) {
		appt := &entity.Appointment{
			Description:         "ขอปรึกษา",
			AdvisorUserID:       3,
			StudentUserID:       9,
			AppointmentStatusID: approval.StatusPending,
			AppointmentStatus:   entity.AppointmentStatus{StatusCode: "PENDING"},
		}
		appt.ID = 1

		repo := newFakeRepo(appt)
		svc := approval.NewAppointmentService(repo)

		updated, err := svc.ApproveAppointment(1, 3, "advisor", "อนุมัติแล้ว")

		Expect(err).To(BeNil())
		Expect(updated).ToNot(BeNil())
		Expect(updated.AppointmentStatusID).To(Equal(approval.StatusApproved))
		Expect(updated.Description).To(Equal("อนุมัติแล้ว"))

		// ✅ UpdateFields ถูกเรียก 2 รอบ: รอบแรกอัปเดต status, รอบสองอัปเดต status_history_id
		Expect(repo.updateCalls).ToNot(BeEmpty())
		Expect(repo.updateCalls[0]["appointment_status_id"]).To(Equal(approval.StatusApproved))

		Expect(repo.lastUpsertID).To(Equal(uint(1)))
		Expect(repo.lastUpsertStatus).To(Equal(approval.StatusApproved))

		Expect(repo.lastHistory).ToNot(BeNil())
		Expect(repo.lastHistory.FromStatusID).To(Equal(approval.StatusPending))
		Expect(repo.lastHistory.ToStatusID).To(Equal(approval.StatusApproved))
		Expect(repo.lastHistory.ActionID).To(Equal(approval.ActionApprove))
		Expect(repo.lastHistory.ChangedByUserID).To(Equal(uint(3)))
		Expect(repo.lastHistory.Reason).To(Equal("อนุมัติแล้ว"))
		Expect(repo.lastHistory.ChangedAt).ToNot(BeZero())
	})

	t.Run("Case 2: Error - role is not advisor", func(t *testing.T) {
		appt := &entity.Appointment{
			AdvisorUserID:       3,
			AppointmentStatusID: approval.StatusPending,
		}
		appt.ID = 1

		repo := newFakeRepo(appt)
		svc := approval.NewAppointmentService(repo)

		updated, err := svc.ApproveAppointment(1, 3, "student", "x")

		Expect(updated).To(BeNil())
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("only advisor can approve appointment"))
	})

	t.Run("Case 3: Error - actor is not the advisor of this appointment", func(t *testing.T) {
		appt := &entity.Appointment{
			AdvisorUserID:       3,
			AppointmentStatusID: approval.StatusPending,
		}
		appt.ID = 1

		repo := newFakeRepo(appt)
		svc := approval.NewAppointmentService(repo)

		updated, err := svc.ApproveAppointment(1, 99, "advisor", "x")

		Expect(updated).To(BeNil())
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("you are not the advisor of this appointment"))
	})

	t.Run("Case 4: Error - appointment is not pending", func(t *testing.T) {
		appt := &entity.Appointment{
			AdvisorUserID:       3,
			AppointmentStatusID: approval.StatusApproved,
		}
		appt.ID = 1

		repo := newFakeRepo(appt)
		svc := approval.NewAppointmentService(repo)

		updated, err := svc.ApproveAppointment(1, 3, "advisor", "x")

		Expect(updated).To(BeNil())
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("appointment is not in pending status"))
	})

	t.Run("Case 5: Error - repo.GetByID fails", func(t *testing.T) {
		repo := newFakeRepo(nil)
		repo.getErr = errors.New("db down")
		svc := approval.NewAppointmentService(repo)

		updated, err := svc.ApproveAppointment(1, 3, "advisor", "x")

		Expect(updated).To(BeNil())
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("db down"))
	})

	t.Run("Case 6: Error - repo.UpdateFields fails", func(t *testing.T) {
		appt := &entity.Appointment{
			AdvisorUserID:       3,
			AppointmentStatusID: approval.StatusPending,
		}
		appt.ID = 1

		repo := newFakeRepo(appt)
		repo.updateFieldsErr = errors.New("update failed")
		svc := approval.NewAppointmentService(repo)

		updated, err := svc.ApproveAppointment(1, 3, "advisor", "x")

		Expect(updated).To(BeNil())
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("update failed"))
	})

	t.Run("Case 7: Error - repo.UpsertAppointmentState fails", func(t *testing.T) {
		appt := &entity.Appointment{
			AdvisorUserID:       3,
			AppointmentStatusID: approval.StatusPending,
		}
		appt.ID = 1

		repo := newFakeRepo(appt)
		repo.upsertStateErr = errors.New("state failed")
		svc := approval.NewAppointmentService(repo)

		updated, err := svc.ApproveAppointment(1, 3, "advisor", "x")

		Expect(updated).To(BeNil())
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("state failed"))
	})

	t.Run("Case 8: Error - repo.CreateStatusHistory fails", func(t *testing.T) {
		appt := &entity.Appointment{
			AdvisorUserID:       3,
			AppointmentStatusID: approval.StatusPending,
		}
		appt.ID = 1

		repo := newFakeRepo(appt)
		repo.createHistoryErr = errors.New("history failed")
		svc := approval.NewAppointmentService(repo)

		updated, err := svc.ApproveAppointment(1, 3, "advisor", "x")

		Expect(updated).To(BeNil())
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("history failed"))
	})
}

// --------------------
// Tests: ProposeNewTime (Reschedule)
// --------------------
func TestProposeNewTime(t *testing.T) {
	RegisterTestingT(t)

	t.Run("Case 1: Success - advisor proposes new time for pending appointment", func(t *testing.T) {
		appt := &entity.Appointment{
			Description:         "ขอปรึกษา",
			AdvisorUserID:       3,
			AppointmentStatusID: approval.StatusPending,
			AppointmentStatus:   entity.AppointmentStatus{StatusCode: "PENDING"},
		}
		appt.ID = 2

		repo := newFakeRepo(appt)
		svc := approval.NewAppointmentService(repo)

		updated, err := svc.ProposeNewTime(2, 3, "ADVISOR", "เสนอเวลาใหม่")

		Expect(err).To(BeNil())
		Expect(updated).ToNot(BeNil())
		Expect(updated.AppointmentStatusID).To(Equal(approval.StatusReschedule))
		Expect(updated.Description).To(Equal("เสนอเวลาใหม่"))

		// ✅ UpdateFields ถูกเรียก 2 รอบเหมือนกัน
		Expect(repo.updateCalls).ToNot(BeEmpty())
		Expect(repo.updateCalls[0]["appointment_status_id"]).To(Equal(approval.StatusReschedule))

		Expect(repo.lastUpsertStatus).To(Equal(approval.StatusReschedule))

		Expect(repo.lastHistory).ToNot(BeNil())
		Expect(repo.lastHistory.FromStatusID).To(Equal(approval.StatusPending))
		Expect(repo.lastHistory.ToStatusID).To(Equal(approval.StatusReschedule))
		Expect(repo.lastHistory.ActionID).To(Equal(approval.ActionReschedule))
		Expect(repo.lastHistory.Reason).To(Equal("เสนอเวลาใหม่"))
		Expect(time.Since(repo.lastHistory.ChangedAt)).To(BeNumerically("<", 5*time.Second))
	})

	t.Run("Case 2: Error - role is not advisor", func(t *testing.T) {
		appt := &entity.Appointment{
			AdvisorUserID:       3,
			AppointmentStatusID: approval.StatusPending,
		}
		appt.ID = 2

		repo := newFakeRepo(appt)
		svc := approval.NewAppointmentService(repo)

		updated, err := svc.ProposeNewTime(2, 3, "student", "x")

		Expect(updated).To(BeNil())
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("only teacher/advisor can propose new time"))
	})

	t.Run("Case 3: Error - appointment is not pending", func(t *testing.T) {
		appt := &entity.Appointment{
			AdvisorUserID:       3,
			AppointmentStatusID: approval.StatusApproved,
		}
		appt.ID = 2

		repo := newFakeRepo(appt)
		svc := approval.NewAppointmentService(repo)

		updated, err := svc.ProposeNewTime(2, 3, "advisor", "x")

		Expect(updated).To(BeNil())
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("appointment is not in pending status"))
	})
}
