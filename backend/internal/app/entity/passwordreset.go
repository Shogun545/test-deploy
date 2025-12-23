package entity

import (
	"gorm.io/gorm"
	"time"
)

type PasswordReset struct {
	gorm.Model
	UserID    	uint      `json:"user_id"`
    User      	*User     `json:"user" gorm:"foreignKey:UserID"`
    ResetToken 	string    `json:"reset_token" gorm:"type:varchar(255);index"`
    ExpiresAt 	time.Time `json:"expires_at"` //แบบว่าตั้งเวลาให้หมดอายุ เวลาเราส่งไปพวก email 5-10 นาทีอะไรแบบนี้
    Used      	bool      `json:"used"` //ซ้ำหรือไม่ซ้ำกับที่เคยใช้
}