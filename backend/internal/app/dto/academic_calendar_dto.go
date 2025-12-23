package dto

type EventRequest struct {
	Title     string `json:"title"`
	Type      string `json:"type"`
	
	// ✅ เปลี่ยนจาก Date เป็น StartDate, EndDate
	StartDate string `json:"start_date"` 
	EndDate   string `json:"end_date"`

	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	
	AdminID   uint   `json:"-"` 
}

type HolidayResponse struct {
	ID          uint     `json:"id"`
	
	// ✅ ส่งกลับเป็น Start/End
	StartDate   string   `json:"start_date"`
	EndDate     string   `json:"end_date"`
	
	LocalName   string   `json:"localName"`
	Name        string   `json:"name"`
	CountryCode string   `json:"countryCode"`
	Types       []string `json:"types"`
	
	StartTime   string   `json:"start_time"`
	EndTime     string   `json:"end_time"`
}