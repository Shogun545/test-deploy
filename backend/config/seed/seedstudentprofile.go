package seed

import (
    "backend/internal/app/entity"
    "gorm.io/gorm"
    "log"
)

func SeedStudentProfile(db *gorm.DB) error {

    // ▸ รายชื่อนักศึกษา 10 คนที่เพิ่มมาใหม่
    studentSutIDs := []string{
        "B6608019", "B6608020", "B6608021", "B6608022", "B6608023", "B6608024",
        "B6608025", "B6608026", "B6608027", "B6608028", "B6608029",
    }

    // ▸ หา advisor ก่อน  (T6608019)
    var advisorUser entity.User
    if err := db.Where("sut_id = ?", "T6608019").First(&advisorUser).Error; err != nil {
        return err
    }

    var advisorProfile entity.AdvisorProfile
    if err := db.Where("user_id = ?", advisorUser.ID).First(&advisorProfile).Error; err != nil {
        return err
    }

    // ▸ loop นักศึกษาทุกคนเพื่อสร้าง profile + assign advisor_profile_id
    for _, sutID := range studentSutIDs {

        // หา user นักศึกษา
        var user entity.User
        if err := db.Where("sut_id = ?", sutID).First(&user).Error; err != nil {
            log.Printf("⚠ หา user ไม่เจอ: %s (skip)", sutID)
            continue
        }

        // Create StudentProfile
        sp := entity.StudentProfile{
            UserID:           user.ID,
            YearOfStudy:      1,
            AdvisorProfileID: &advisorProfile.ID, // ← ผูกอาจารย์ที่ปรึกษา
        }

        if err := db.FirstOrCreate(&sp, entity.StudentProfile{
            UserID: user.ID,
        }).Error; err != nil {
            return err
        }

        // สร้าง Academic Record ให้ด้วย (optional)
        record := entity.StudentAcademicRecord{
            StudentProfileID: sp.ID,
            AcademicYear:     "2568",
            Semester:         1,
            TermGPA:          0,
            CumulativeGPA:    0,
            AcademicStatus:   "ปกติ",
        }

        db.FirstOrCreate(&record, entity.StudentAcademicRecord{
            StudentProfileID: sp.ID,
            AcademicYear:     "2568",
            Semester:         1,
        })
    }

    log.Println("Seeded 10 students with advisor T6608019")
    return nil
}
