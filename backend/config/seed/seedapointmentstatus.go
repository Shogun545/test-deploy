package seed

import (
    "log"
    "gorm.io/gorm"
    "backend/internal/app/entity"
)

func SeedAppointmentStatus(db *gorm.DB) {
    statuses := []entity.AppointmentStatus{
        {
            Model:        gorm.Model{ID: 1},
            StatusCode:   "PENDING",
            StatusName:   "รอพิจารณา",
            IsTerminal:   false,
            DisplayOrder: 1,
        },
        {
            Model:        gorm.Model{ID: 2},
            StatusCode:   "APPROVED",
            StatusName:   "อนุมัติแล้ว",
            IsTerminal:   true,
            DisplayOrder: 2,
        },
        {
            Model:        gorm.Model{ID: 3},
            StatusCode:   "RESCHEDULE",
            StatusName:   "เสนอเวลาใหม่",
            IsTerminal:   false,
            DisplayOrder: 3,
        },
    }

    for _, status := range statuses {

        cond := entity.AppointmentStatus{Model: gorm.Model{ID: status.ID}}

        if err := db.FirstOrCreate(&status, cond).Error; err != nil {
            log.Printf("Failed to seed AppointmentStatus %s: %v", status.StatusCode, err)
        } else {
            log.Printf("Seeded AppointmentStatus: %s", status.StatusCode)
        }
    }
}
