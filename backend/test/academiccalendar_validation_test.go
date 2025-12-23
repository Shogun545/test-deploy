package test_test 

import (
	"testing"
	"time"
	"backend/internal/app/entity" // ✅ Import entity เข้ามาเพื่อ Test
	. "github.com/onsi/gomega"
)

func TestAcademicCalendarEntityValidation(t *testing.T) {
	RegisterTestingT(t)

	// ตั้งค่าเวลาจำลอง
	baseTime := time.Date(2025, 12, 20, 9, 0, 0, 0, time.Local)

	// --- กรณีที่ 1: ข้อมูลถูกต้อง (Success) ---
	t.Run("Case 1: Valid Entity", func(t *testing.T) {
		event := entity.AcademicCalendar{
			EventName:     "สอบกลางภาค",
			EventType:     "exam",
			StartDateTime: baseTime,
			EndDateTime:   baseTime.Add(3 * time.Hour), // จบหลังเริ่ม 3 ชม.
		}

		// ✅ เรียกฟังก์ชัน Validate ที่อยู่ใน Entity   
		err := event.Validate()
		Expect(err).To(BeNil())
	})

	// --- กรณีที่ 2: ชื่อกิจกรรมว่าง (Required Field) ---
	t.Run("Case 2: EventName is empty", func(t *testing.T) {
		event := entity.AcademicCalendar{
			EventName:     "", // ❌ ว่าง
			EventType:     "exam",
			StartDateTime: baseTime,
			EndDateTime:   baseTime.Add(1 * time.Hour),
		}

		err := event.Validate()

		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(ContainSubstring("กรุณาระบุชื่อกิจกรรม"))
	})

	// // --- กรณีที่ 3: ประเภทกิจกรรมผิด (Enum) ---
	t.Run("Case 3: Invalid EventType", func(t *testing.T) {
		event := entity.AcademicCalendar{
			EventName:     "ปาร์ตี้",
			EventType:     "party", // ❌ ไม่อยู่ใน exam, activity, holiday
			StartDateTime: baseTime,
			EndDateTime:   baseTime.Add(1 * time.Hour),
		}

		err := event.Validate()

		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(ContainSubstring("ประเภทกิจกรรมไม่ถูกต้อง"))
	})

	// // --- กรณีที่ 4: เวลาผิด Logic (จบก่อนเริ่ม) ---
	t.Run("Case 4: EndTime before StartTime", func(t *testing.T) {
		event := entity.AcademicCalendar{
			EventName:     "สอบ",
			EventType:     "exam",
			StartDateTime: baseTime,
			EndDateTime:   baseTime.Add(-1 * time.Hour), // ❌ จบก่อนเริ่ม (ย้อนหลัง 1 ชม.)
		}

		err := event.Validate()

		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("วันเวลาสิ้นสุด ต้องอยู่หลังจากวันเวลาเริ่มต้น"))
	})
}