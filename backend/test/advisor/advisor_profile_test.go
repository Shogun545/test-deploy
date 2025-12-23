package test

import (
	"testing"

	"backend/internal/app/entity"
	. "github.com/onsi/gomega"
)

func TestAdvisorProfile_Valid(t *testing.T) {
	g := NewWithT(t)

	// ตัวอย่าง AdvisorProfile ถูกต้อง
	advisor := &entity.AdvisorProfile{
		User:        &entity.User{Phone: "0629931709"}, // Phone ต้องกรอก
		OfficeRoom:  "A101",                             // ห้องทำงาน
		Specialties: "Computer Architecture",            // ความเชี่ยวชาญ
		IsActive:    true,                               // Active หรือไม่
	}

	// คาดหวังว่า Validate จะผ่าน ไม่มี error
	g.Expect(advisor.Validate()).To(BeNil())
}

func TestAdvisorProfile_EmptyPhone(t *testing.T) {
	g := NewWithT(t)

	advisor := &entity.AdvisorProfile{
		User:        &entity.User{Phone: ""},           // เบอร์ว่าง
		OfficeRoom:  "A101",
		Specialties: "Computer Architecture",
	}

	// คาดหวังว่า Validate จะ error
	err := advisor.Validate()
	g.Expect(err).ToNot(BeNil())
	g.Expect(err.Error()).To(Equal("phone is required"))
}

func TestAdvisorProfile_InvalidPhone(t *testing.T) {
	g := NewWithT(t)

	advisor := &entity.AdvisorProfile{
		User:        &entity.User{Phone: "12345"},      // เบอร์ไม่ครบ 10 หลัก
		OfficeRoom:  "A101",
		Specialties: "Computer Architecture",
	}

	err := advisor.Validate()
	g.Expect(err).ToNot(BeNil())
	g.Expect(err.Error()).To(Equal("phone must be 10 digits"))
}

func TestAdvisorProfile_EmptyOfficeRoom(t *testing.T) {
	g := NewWithT(t)

	advisor := &entity.AdvisorProfile{
		User:        &entity.User{Phone: "0629931709"},
		OfficeRoom:  "",                                 // ห้องว่าง
		Specialties: "Computer Architecture",
	}

	err := advisor.Validate()
	g.Expect(err).ToNot(BeNil())
	g.Expect(err.Error()).To(Equal("officeRoom is required"))
}

func TestAdvisorProfile_EmptySpecialties(t *testing.T) {
	g := NewWithT(t)

	advisor := &entity.AdvisorProfile{
		User:        &entity.User{Phone: "0629931709"},
		OfficeRoom:  "A101",
		Specialties: "",                                 // ความเชี่ยวชาญว่าง
	}

	err := advisor.Validate()
	g.Expect(err).ToNot(BeNil())
	g.Expect(err.Error()).To(Equal("specialties is required"))
}
