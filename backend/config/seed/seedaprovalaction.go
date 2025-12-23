package seed

import (
	"backend/internal/app/entity"
	"gorm.io/gorm"
)

func SeedApprovalActions(db *gorm.DB) {
	actions := []entity.ApprovalAction{
		{
			Model:        gorm.Model{ID: 1},
			ActionCode:   "APPROVE",
			ActionName:   "Approve appointment",
			IsActive:     true,
			DisplayOrder: 1,
		},
		{
			Model:        gorm.Model{ID: 2},
			ActionCode:   "RESCHEDULE",
			ActionName:   "Propose new time",
			IsActive:     true,
			DisplayOrder: 2,
		},
	}

	for _, action := range actions {
		var existing entity.ApprovalAction
		err := db.First(&existing, action.ID).Error
		if err != nil {
			db.Create(&action)  // สร้างใหม่เมื่อไม่มี
		}
	}
}
