package seed

import (
    "backend/internal/app/entity"
    "gorm.io/gorm"
)

func SeedAdvisorProfile(db *gorm.DB) error {

    // 1) เลือก advisor จาก sut_id (ต้องตรงกับ user ที่ login ได้)
    var user entity.User
    if err := db.Where("sut_id = ?", "T6608019").First(&user).Error; err != nil {
        return err
    }

    // 2) เตรียมข้อมูล AdvisorProfile
    profile := entity.AdvisorProfile{
        UserID:        user.ID,
        Specialties:   "ที่ปรึกษาด้านวิศวกรรมคอมพิวเตอร์", // ความเชี่ยวชาญ
        IsActive:      true,                                   //เปิดให้ปรึกษา
        OfficeRoom:    "B-201",                                //ห้องทำงาน
    }

    // 3) บันทึกแบบกัน duplicate
    if err := db.FirstOrCreate(&profile, entity.AdvisorProfile{
        UserID: user.ID,
    }).Error; err != nil {
        return err
    }

    return nil
}
