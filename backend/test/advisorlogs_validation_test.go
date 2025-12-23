package test

import (
	"testing"

	"backend/internal/app/entity"
	"github.com/asaskevich/govalidator"
	. "github.com/onsi/gomega"
)

func TestAdvisorLogValidation(t *testing.T) {
	RegisterTestingT(t)

	// ✅ Positive 1 case
	t.Run("UT-01: valid advisor log", func(t *testing.T) {
		log := entity.AdvisorLog{
			AppointmentID:  1,
			Title:          "สรุปการปรึกษาครั้งที่ 1",
			Body:           "นักศึกษาทำได้ดี มีงานที่ต้องแก้ไขต่อ",
			Status:         "PendingReport",
			RequiresReport: true,
		}

		ok, err := govalidator.ValidateStruct(log)
		Expect(ok).To(BeTrue())
		Expect(err).To(BeNil())
	})

	// ❌ Negative 3 cases
	t.Run("UT-02: appointmentId is required", func(t *testing.T) {
		log := entity.AdvisorLog{
			AppointmentID: 0,
			Title:         "Topic",
			Body:          "Detail",
			Status:        "Completed",
		}

		ok, err := govalidator.ValidateStruct(log)
		Expect(ok).To(BeFalse())
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(ContainSubstring("appointmentId is required"))
	})

	t.Run("UT-03: title is required", func(t *testing.T) {
		log := entity.AdvisorLog{
			AppointmentID: 1,
			Title:         "",
			Body:          "Detail",
			Status:        "Completed",
		}

		ok, err := govalidator.ValidateStruct(log)
		Expect(ok).To(BeFalse())
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(ContainSubstring("title is required"))
	})

	t.Run("UT-04: body is required", func(t *testing.T) {
		log := entity.AdvisorLog{
			AppointmentID: 1,
			Title:         "Topic",
			Body:          "",
			Status:        "Completed",
		}

		ok, err := govalidator.ValidateStruct(log)
		Expect(ok).To(BeFalse())
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(ContainSubstring("body is required"))
	})
}
