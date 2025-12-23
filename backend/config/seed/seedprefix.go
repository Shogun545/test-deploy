package seed

import (
    "log"

    "gorm.io/gorm"
    "backend/internal/app/entity"
)

func SeedPrefix(db *gorm.DB) {

    prefixes := []entity.Prefix{
        {Prefix: "นาย"},
        {Prefix: "นาง"},
        {Prefix: "นางสาว"},
        {Prefix: "ดร."},
		{Prefix: "ผศ."},
		{Prefix: "ผศ.ดร."},
		{Prefix: "รศ."},
		{Prefix: "รศ.ดร."},
		{Prefix: "ศ."},
		{Prefix: "ศ.ดร."},
    }

    for _, p := range prefixes {
        var existing entity.Prefix
        err := db.Where("prefix = ?", p.Prefix).First(&existing).Error

        if err == gorm.ErrRecordNotFound {
            db.Create(&p)
        } else if err == nil {
            log.Printf("Prefix '%s' already exists, skipping...", p.Prefix)
        }
    }

    log.Println("Prefix seed completed.")
}
