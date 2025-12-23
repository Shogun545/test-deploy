package auth_test

import (
    "testing"
    "backend/internal/app/entity"
    . "github.com/onsi/gomega"
)

// Test Login struct (SutId และ Password ไม่ว่าง)
func TestLogin_Valid(t *testing.T) {
    g := NewWithT(t)

    // กรณี login ถูกต้องครบถ้วน
    login := entity.Login{
        SutId:    "B6608019",
        Password: "123456",
    }

    // คาดหวังว่า Validate ผ่าน (ไม่เกิด error)
    g.Expect(login.Validate()).To(BeNil())
}

func TestLogin_EmptySutId(t *testing.T) {
    g := NewWithT(t)

    // SutId เป็นค่าว่าง
    login := entity.Login{
        SutId:    "",
        Password: "123456",
    }

    err := login.Validate()

    // คาดหวังว่าเกิด error และข้อความตรง
    g.Expect(err).ToNot(BeNil())
    g.Expect(err.Error()).To(Equal("sut_id is required"))
}

func TestLogin_EmptyPassword(t *testing.T) {
    g := NewWithT(t)

    // Password เป็นค่าว่าง
    login := entity.Login{
        SutId:    "B6608019",
        Password: "",
    }

    err := login.Validate()

    // คาดหวังว่าเกิด error และข้อความตรง
    g.Expect(err).ToNot(BeNil())
    g.Expect(err.Error()).To(Equal("password is required"))
}

//  Test User sut_id และ Password strength
func TestUserLoginValidation(t *testing.T) {
    g := NewWithT(t)

    // กรณีปกติ ข้อมูลครบถ้วน
    user := entity.User{
        SutId:        "B6608019",
        PasswordHash: "Strong@123",
        Role:         &entity.Role{Role: "student"},
    }

    // ตรวจสอบ sut_id ขึ้นต้นถูกต้อง
    g.Expect(user.ValidateSutIdByRole()).To(BeNil())

    // ตรวจสอบความแข็งแรงของ password
    g.Expect(user.ValidatePasswordStrength()).To(BeNil())
}

func TestUser_InvalidSutId(t *testing.T) {
    g := NewWithT(t)

    // sut_id ผิด (student แต่ขึ้นต้นไม่ใช่ B)
    user := entity.User{
        SutId:        "T6608019",
        PasswordHash: "Strong@123",
        Role:         &entity.Role{Role: "student"},
    }

    err := user.ValidateSutIdByRole()
    g.Expect(err).ToNot(BeNil())
    g.Expect(err.Error()).To(Equal("student sut_id must start with B"))
}

func TestUser_WeakPassword(t *testing.T) {
    g := NewWithT(t)

    // password อ่อน (ไม่มีตัวใหญ่)
    user := entity.User{
        SutId:        "B6608019",
        PasswordHash: "weak@123",
        Role:         &entity.Role{Role: "student"},
    }

    err := user.ValidatePasswordStrength()
    g.Expect(err).ToNot(BeNil())
}

func TestPasswordStrength_AllCases(t *testing.T) {
    g := NewWithT(t)

    // 1. Password ถูกต้องครบทุกข้อ
    strongPw := entity.User{PasswordHash: "Strong@123"}
    g.Expect(strongPw.ValidatePasswordStrength()).To(BeNil())

    // 2. Password สั้นเกินไป
    tooShort := entity.User{PasswordHash: "S@1a"}
    g.Expect(tooShort.ValidatePasswordStrength()).ToNot(BeNil())
    g.Expect(tooShort.ValidatePasswordStrength().Error()).To(Equal("password must be at least 8 characters long"))

    // 3. ไม่มีตัวพิมพ์ใหญ่
    noUpper := entity.User{PasswordHash: "weak@1234"}
    g.Expect(noUpper.ValidatePasswordStrength()).ToNot(BeNil())
    g.Expect(noUpper.ValidatePasswordStrength().Error()).To(Equal("password must contain at least one uppercase letter"))

    // 4. ไม่มีตัวพิมพ์เล็ก
    noLower := entity.User{PasswordHash: "WEAK@1234"}
    g.Expect(noLower.ValidatePasswordStrength()).ToNot(BeNil())
    g.Expect(noLower.ValidatePasswordStrength().Error()).To(Equal("password must contain at least one lowercase letter"))

    // 5. ไม่มีตัวเลข
    noNumber := entity.User{PasswordHash: "Weak@pass"}
    g.Expect(noNumber.ValidatePasswordStrength()).ToNot(BeNil())
    g.Expect(noNumber.ValidatePasswordStrength().Error()).To(Equal("password must contain at least one number"))

    // 6. ไม่มีสัญลักษณ์พิเศษ
    noSpecial := entity.User{PasswordHash: "Weak1234"}
    g.Expect(noSpecial.ValidatePasswordStrength()).ToNot(BeNil())
    g.Expect(noSpecial.ValidatePasswordStrength().Error()).To(Equal("password must contain at least one special character"))
}