package entity

import "github.com/asaskevich/govalidator"

// Login คือ struct ที่แทน "ฟอร์ม login" ไม่ใช่ตารางใน DB ใช้เฉพาะรับข้อมูลจาก frontend ตอน login
type Login struct {
	SutId    string `json:"sut_id" valid:"required~sut_id is required"`
	Password string `json:"password" valid:"required~password is required"`
}

// Validate ใช้ตรวจสอบว่า input จาก frontend ถูกต้องหรือไม่ ไม่เช็ค role,prefix,เช็คแค่ว่าไม่ว่าง
func (l *Login) Validate() error {
	_, err := govalidator.ValidateStruct(l)
	return err
}
