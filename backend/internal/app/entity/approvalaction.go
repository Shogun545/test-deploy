package entity
import ("gorm.io/gorm")
type ApprovalAction struct {
	gorm.Model
    
    ActionCode   string `gorm:"type:varchar(100)" json:"action_code"`
    ActionName   string `gorm:"type:varchar(100)" json:"action_name"`

    IsActive     bool     `gorm:"type:boolean" json:"is_active"`
    DisplayOrder int        `gorm:"type:int" json:"display_order"`
}
