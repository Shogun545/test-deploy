package test

import (
    "testing"

    "backend/internal/app/entity"
    . "github.com/onsi/gomega"
)


// TestStudentProfile_Valid: ทดสอบกรณีข้อมูลนักศึกษาถูกต้องครบถ้วน
func TestStudentProfile_Valid(t *testing.T) {
    g := NewWithT(t)

    user := entity.User{
        Role:     &entity.Role{Role: "student"}, // กำหนด Role เป็น student
        Birthday: "2005-06-30",                 // วันเกิดในรูปแบบ string (YYYY-MM-DD)
        Phone:    "0629931709",                 // เบอร์โทรศัพท์ครบ 10 หลัก
    }

    g.Expect(user.ValidateStudentProfile()).To(BeNil())
}

// TestStudentProfile_EmptyBirthday: ทดสอบกรณีไม่ระบุวันเกิด
func TestStudentProfile_EmptyBirthday(t *testing.T) {
    g := NewWithT(t)

    user := entity.User{
        Role:  &entity.Role{Role: "student"},
        Phone: "0629931709",
        // ไม่ได้กำหนดค่า Birthday
    }

    err := user.ValidateStudentProfile()
    // ต้องเกิด Error เพราะ Birthday เป็นค่าว่าง
    g.Expect(err).ToNot(BeNil())
    g.Expect(err.Error()).To(Equal("birthday is required"))
}

// TestStudentProfile_EmptyPhone: ทดสอบกรณีไม่ระบุเบอร์โทรศัพท์
func TestStudentProfile_EmptyPhone(t *testing.T) {
    g := NewWithT(t)

    user := entity.User{
        Role:     &entity.Role{Role: "student"},
        Birthday: "2005-06-30",
        // ไม่ได้กำหนดค่า Phone
    }

    err := user.ValidateStudentProfile()
    // ต้องเกิด Error เพราะ Phone เป็นค่าว่าง
    g.Expect(err).ToNot(BeNil())
    g.Expect(err.Error()).To(Equal("phone is required"))
}

// TestStudentProfile_InvalidPhone: ทดสอบกรณีเบอร์โทรศัพท์ไม่ครบ 10 หลัก
func TestStudentProfile_InvalidPhone(t *testing.T) {
    g := NewWithT(t)

    user := entity.User{
        Role:     &entity.Role{Role: "student"},
        Birthday: "2005-06-30",
        Phone:    "12345", // ใส่เบอร์สั้นเกินไป
    }

    err := user.ValidateStudentProfile()
    // ต้องเกิด Error และแจ้งเตือนเรื่องจำนวนหลักของเบอร์โทร
    g.Expect(err).ToNot(BeNil())
    g.Expect(err.Error()).To(Equal("phone must be 10 digits"))
}