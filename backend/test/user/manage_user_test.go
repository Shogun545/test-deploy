package test

import (
	"testing"
	"backend/internal/app/entity"
	. "github.com/onsi/gomega"
)

// กรณีทดสอบ: sut_id ของ student ถูกต้อง
func TestStudentSutId_Valid(t *testing.T) {
	g := NewWithT(t)

	user := entity.User{
		SutId: "B6608019",                 // ยาว 8 ตัวและขึ้นต้นด้วย B
		Role:  &entity.Role{Role: "student"}, // role เป็น student
	}

	err := user.ValidateSutIdByRole()
	g.Expect(err).To(BeNil()) // คาดหวังว่าไม่เกิด error
}

// กรณีทดสอบ: sut_id ของ student ผิดตัวขึ้นต้น
func TestStudentSutId_InvalidPrefix(t *testing.T) {
	g := NewWithT(t)

	user := entity.User{
		SutId: "T6608019", // ขึ้นต้น T แต่ role เป็น student
		Role:  &entity.Role{Role: "student"},
	}

	err := user.ValidateSutIdByRole()
	g.Expect(err).ToNot(BeNil())                  // ต้องเกิด error
	g.Expect(err.Error()).To(Equal("student sut_id must start with B")) // ตรวจข้อความ
}

// กรณีทดสอบ: sut_id ของ student ผิดความยาว
func TestStudentSutId_InvalidLength(t *testing.T) {
	g := NewWithT(t)

	user := entity.User{
		SutId: "B660801", // 7 ตัว
		Role:  &entity.Role{Role: "student"},
	}

	err := user.ValidateSutIdByRole()
	g.Expect(err).ToNot(BeNil())
	g.Expect(err.Error()).To(Equal("sut_id must be 8 characters"))
}

// กรณีทดสอบ: sut_id ของ advisor ผิดตัวขึ้นต้น
func TestAdvisorSutId_InvalidPrefix(t *testing.T) {
	g := NewWithT(t)

	user := entity.User{
		SutId: "B6608019",                 // ขึ้นต้น B แต่ role เป็น advisor
		Role:  &entity.Role{Role: "advisor"},
	}

	err := user.ValidateSutIdByRole()
	g.Expect(err).ToNot(BeNil())
	g.Expect(err.Error()).To(Equal("advisor sut_id must start with T"))
}

// กรณีทดสอบ: sut_id ของ admin ถูกต้อง
func TestAdminSutId_Valid(t *testing.T) {
	g := NewWithT(t)

	user := entity.User{
		SutId: "A0010001",                  // ขึ้นต้น A ยาว 8 ตัว
		Role:  &entity.Role{Role: "admin"}, // role เป็น admin
	}

	err := user.ValidateSutIdByRole()
	g.Expect(err).To(BeNil()) // ต้องผ่าน validation
}

// กรณีทดสอบ: user active
func TestManageUser_Active(t *testing.T) {
	g := NewWithT(t)

	user := &entity.User{
		Active: true,
	}

	// คาดหวังว่า Validate ผ่าน
	g.Expect(user.ValidateActiveStatus()).To(BeNil())
	g.Expect(user.Active).To(BeTrue())
}

// กรณีทดสอบ: user inactive
func TestManageUser_Inactive(t *testing.T) {
	g := NewWithT(t)

	user := &entity.User{
		Active: false,
	}

	// คาดหวังว่า Validate ผ่าน
	g.Expect(user.ValidateActiveStatus()).To(BeNil())
	g.Expect(user.Active).To(BeFalse())
}

