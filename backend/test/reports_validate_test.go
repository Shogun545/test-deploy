package test

import (
	"testing"

	"backend/internal/app/entity"
	"github.com/asaskevich/govalidator"
	. "github.com/onsi/gomega"
)

func TestReportValidation(t *testing.T) {
	RegisterTestingT(t)

	//test all correct
	t.Run("valid: all required fields", func(t *testing.T) {
		report := entity.Report{
			Description:    "ไม่สามารถเช้าถึงหน้า UI ได้ไม่ทราบสาเหตุ",
			ReportByID:     1,
			ReportStatusID: 1,
			ReportTopicID:  1,
		}

		ok, err := govalidator.ValidateStruct(report)
		Expect(ok).To(BeTrue())
		Expect(err).To(BeNil())
	})
	//test description null
	t.Run("invalid: description is required", func(t *testing.T) {
		report := entity.Report{
			Description:    "",
			ReportByID:     1,
			ReportStatusID: 1,
			ReportTopicID:  1,
		}

		ok, err := govalidator.ValidateStruct(report)
		Expect(ok).To(BeFalse())
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(ContainSubstring("description is required"))
	})
	//test userid null
	t.Run("invalid: reportById is required", func(t *testing.T) {
		report := entity.Report{
			Description:    "รายงาน",
			ReportByID:     0,
			ReportStatusID: 1,
			ReportTopicID:  1,
		}

		ok, err := govalidator.ValidateStruct(report)
		Expect(ok).To(BeFalse())
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(ContainSubstring("reportById is required"))
	})
	//test reportstatus null
	t.Run("invalid: reportStatusId is required", func(t *testing.T) {
		report := entity.Report{
			Description:    "รายงาน",
			ReportByID:     1,
			ReportStatusID: 0,
			ReportTopicID:  1,
		}

		ok, err := govalidator.ValidateStruct(report)
		Expect(ok).To(BeFalse())
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(ContainSubstring("reportStatusId is required"))
	})
	//test reporttopic null
	t.Run("invalid: reportTopicId is required", func(t *testing.T) {
		report := entity.Report{
			Description:    "รายงาน",
			ReportByID:     1,
			ReportStatusID: 1,
			ReportTopicID:  0,
		}

		ok, err := govalidator.ValidateStruct(report)
		Expect(ok).To(BeFalse())
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(ContainSubstring("reportTopicId is required"))
	})
}
