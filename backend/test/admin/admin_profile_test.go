package test

import (
	"testing"
	"backend/internal/app/entity"
	. "github.com/onsi/gomega" 
)

// TestAdminProfile_Valid: ทดสอบกรณีข้อมูลทุกอย่างถูกต้อง
func TestAdminProfile_Valid(t *testing.T) {
	g := NewWithT(t) // สร้าง gomega matcher

	admin := &entity.User{
		Email:  "admin@sut.ac.th",
		Phone:  "0629931709",
		Notes:  "ดูแลระบบทั้งหมด",
		Active: true,
	}
	g.Expect(admin.ValidateAdminProfile()).To(BeNil())
}

// TestAdminProfile_EmptyEmail: ทดสอบกรณีไม่ได้กรอก Email
func TestAdminProfile_EmptyEmail(t *testing.T) {
	g := NewWithT(t)

	admin := &entity.User{
		Email:  "", // จำลองค่าว่าง
		Phone:  "0629931709",
		Notes:  "ดูแลระบบทั้งหมด",
		Active: true,
	}

	err := admin.ValidateAdminProfile()
	g.Expect(err).ToNot(BeNil())
	g.Expect(err.Error()).To(Equal("email is required"))
}

// TestAdminProfile_InvalidEmail: ทดสอบกรณีรูปแบบ Email ไม่ถูกต้อง (ไม่มี @)
func TestAdminProfile_InvalidEmail(t *testing.T) {
	g := NewWithT(t)

	admin := &entity.User{
		Email:  "admin_at_sut.ac.th", // รูปแบบผิด
		Phone:  "0629931709",
		Notes:  "ดูแลระบบทั้งหมด",
		Active: true,
	}

	err := admin.ValidateAdminProfile()
	g.Expect(err).ToNot(BeNil())
	g.Expect(err.Error()).To(Equal("email is invalid"))
}

// TestAdminProfile_EmptyPhone: ทดสอบกรณีไม่ได้กรอกเบอร์โทรศัพท์
func TestAdminProfile_EmptyPhone(t *testing.T) {
	g := NewWithT(t)

	admin := &entity.User{
		Email:  "admin@sut.ac.th",
		Phone:  "", // จำลองค่าว่าง
		Notes:  "ดูแลระบบทั้งหมด",
		Active: true,
	}

	err := admin.ValidateAdminProfile()
	g.Expect(err).ToNot(BeNil())
	g.Expect(err.Error()).To(Equal("phone is required"))
}

// TestAdminProfile_InvalidPhone: ทดสอบกรณีเบอร์โทรศัพท์มีความยาวไม่ถูกต้อง
func TestAdminProfile_InvalidPhone(t *testing.T) {
	g := NewWithT(t)

	admin := &entity.User{
		Email:  "admin@sut.ac.th",
		Phone:  "12345", // ใส่แค่ 5 หลัก (ผิดเงื่อนไขที่ต้องมี 10 หลัก)
		Notes:  "ดูแลระบบทั้งหมด",
		Active: true,
	}

	err := admin.ValidateAdminProfile()
	g.Expect(err).ToNot(BeNil())
	g.Expect(err.Error()).To(Equal("phone must be 10 digits"))
}