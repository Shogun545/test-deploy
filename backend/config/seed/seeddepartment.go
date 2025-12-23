package seed

import (
	"log"

	"backend/internal/app/entity"
	"gorm.io/gorm"
)

func SeedDepartmentData(db *gorm.DB) {
	departments := []entity.Department{
		{Department: "IT Department"},      
	}

	for _, department := range departments {
		if err := db.FirstOrCreate(&department, &entity.Department{Department: department.Department}).Error; err != nil {
			log.Fatalf("failed to seed department %s: %v", department.Department, err)
		}
	}
}