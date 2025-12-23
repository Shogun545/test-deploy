package dto

// User (ที่แสดงใน Report)
type ReportUserDTO struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// Status
type ReportStatusDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// Topic
type ReportTopicDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// Report
type ReportDTO struct {
	ID          uint               `json:"id"`
	Description string             `json:"description"`
	User        *ReportUserDTO     `json:"user"`
	Status      *ReportStatusDTO   `json:"status"`
	Topic       *ReportTopicDTO    `json:"topic"`
}
