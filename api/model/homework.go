package model

import (
	"time"
)

type Homeworks []*Homework

type Homework struct {
	ID              uint `gorm:"primaryKey"`
	Content         string
	Grade           string
	Comments        string
	SubmittedAt     time.Time
	UpdatedAt       time.Time
	StudentFullname string
	StudentID       uint `gorm:"default:null"`
	HomeworkPageID  uint
}

type HomeworkDtos []*HomeworkDto

type HomeworkDto struct {
	ID              uint   `json:"id"`
	Content         string `json:"content"`
	Grade           string `json:"grade"`
	Comments        string `json:"comments"`
	SubmittedAt     string `json:"submitted_at"`
	UpdatedAt       string `json:"updated_at"`
	StudentFullname string `json:"student_fullname"`
	StudentID       uint   `json:"student_id"`
	HomeworkPageID  uint   `json:"homework_page_id"`
}

type HomeworkForm struct {
	Title           string `json:"title" form:"required,alpha_space,max=255"`
	Content         string `json:"content" form:"required,max=255"`
	Grade           string `json:"grade" form:"max=255"`
	Comments        string `json:"comments" form:"max=255"`
	StudentFullname string `json:"student_fullname" form:"alpha_space,max=255"`
	HomeworkPageID  uint   `json:"homework_page_id" form:""`
}

func (hw Homework) ToDto() *HomeworkDto {
	return &HomeworkDto{
		ID:              hw.ID,
		Content:         hw.Content,
		Grade:           hw.Grade,
		Comments:        hw.Comments,
		SubmittedAt:     hw.SubmittedAt.Format(time.RFC3339),
		UpdatedAt:       hw.UpdatedAt.Format(time.RFC3339),
		StudentFullname: hw.StudentFullname,
		StudentID:       hw.StudentID,
		HomeworkPageID:  hw.HomeworkPageID,
	}
}

func (hws Homeworks) ToDto() HomeworkDtos {
	dtos := make([]*HomeworkDto, len(hws))
	for i, hw := range hws {
		dtos[i] = hw.ToDto()
	}
	return dtos
}

func (f *HomeworkForm) ToModel() (*Homework, error) {
	return &Homework{
		Content:         f.Content,
		Grade:           f.Grade,
		Comments:        f.Comments,
		SubmittedAt:     time.Now(),
		StudentFullname: f.StudentFullname,
		HomeworkPageID:  f.HomeworkPageID,
	}, nil
}
