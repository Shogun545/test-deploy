package seed

import (
	"log"

	"gorm.io/gorm"
	"backend/internal/app/entity"
)

func SeedRoleData(db *gorm.DB) {
	// Roles
	roles := []entity.Role{
		{Role: "Admin"},
		{Role: "Student"},
		{Role: "Advisor"},
	}
	for _, role := range roles {
		if err := db.FirstOrCreate(&role, &entity.Role{Role: role.Role}).Error; err != nil {
			log.Fatalf("failed to seed role %s: %v", role.Role, err)
		}
	}
}
