package dto

import "mime/multipart"

type AdvisorLogCreateReq struct {
	AppointmentID  uint                    `form:"appointmentId" binding:"required"`
	Title          string                  `form:"title" binding:"required"`
	Body           string                  `form:"body" binding:"required"`
	RequiresReport bool                    `form:"requiresReport"`
	Status         string                  `form:"status"` 
	Files          []*multipart.FileHeader `form:"files"`  
}

type AdvisorLogUpdateStatusReq struct {
	Status string `json:"status" binding:"required"`
}

type AdvisorLogUpdateReq struct {
	Title          *string                 `form:"title"`
	Body           *string                 `form:"body"`
	RequiresReport *bool                   `form:"requiresReport"`
	Files          []*multipart.FileHeader `form:"files"`
}


type AdvisorLogRespBase struct {
	ID             uint   `json:"id"`
	AppointmentID  uint   `json:"appointmentId"`
	Title          string `json:"title"`
	Body           string `json:"body"`
	Status         string `json:"status"`
	RequiresReport bool   `json:"requiresReport"`
	FileName       string `json:"fileName"`
	FilePath       string `json:"-"` 

	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type AdvisorLogCreateResp struct{ AdvisorLogRespBase }
type AdvisorLogGetResp struct{ AdvisorLogRespBase }
type AdvisorLogUpdateResp struct{ AdvisorLogRespBase }

type AdvisorLogListItemResp struct {
	AdvisorLogRespBase
}