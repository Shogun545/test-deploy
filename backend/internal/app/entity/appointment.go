package entity
import ("gorm.io/gorm"
)
type Appointment struct {
	gorm.Model
    Description string `gorm:"type:text" json:"description"`


    AdvisorUserID       uint                `json:"advisor_user_id"`
    AdvisorUser         User                `gorm:"foreignKey:AdvisorUserID" json:"advisor_user"`
    StudentUserID       uint                `json:"student_user_id"`
    StudentUser         User                `gorm:"foreignKey:StudentUserID" json:"student_user"`
    TopicID             uint                `json:"topic_id"`
    Topic               AppointmentTopic    `gorm:"foreignKey:TopicID"`
    CategoryID          uint                `json:"category_id"`
    Category            AppointmentCategory `gorm:"foreignKey:CategoryID"`
    AppointmentStatusID uint                `json:"appointment_status_id"`
    AppointmentStatus   AppointmentStatus   `gorm:"foreignKey:AppointmentStatusID"`
    AdvisorLog *AdvisorLog `gorm:"foreignKey:AppointmentID" json:"advisor_log"`
}
