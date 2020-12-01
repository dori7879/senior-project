package model

import (
	"time"
)

type HomeworkPages []*HomeworkPage

type HomeworkPage struct {
	ID              uint `gorm:"primaryKey"`
	Title           string
	Content         string
	StudentLink     string `gorm:"unique"`
	TeacherLink     string `gorm:"unique"`
	CourseTitle     string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	OpenedAt        time.Time
	ClosedAt        time.Time
	TeacherFullname string
	TeacherID       uint       `gorm:"default:null"`
	Homeworks       []Homework `gorm:"foreignKey:HomeworkPageID;references:ID"`
}

type HomeworkPageDtos []*HomeworkPageDto

type HomeworkPageDto struct {
	ID              uint   `json:"id"`
	Title           string `json:"title"`
	Content         string `json:"content"`
	StudentLink     string `json:"student_link"`
	TeacherLink     string `json:"teacher_link"`
	CourseTitle     string `json:"course_title"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	OpenedAt        string `json:"opened_at"`
	ClosedAt        string `json:"closed_at"`
	TeacherFullname string `json:"teacher_fullname"`
	TeacherID       uint   `json:"teacher_id"`
}

type HomeworkPageNestedDto struct {
	ID              uint         `json:"id"`
	Title           string       `json:"title"`
	Content         string       `json:"content"`
	StudentLink     string       `json:"student_link"`
	TeacherLink     string       `json:"teacher_link"`
	CourseTitle     string       `json:"course_title"`
	CreatedAt       string       `json:"created_at"`
	UpdatedAt       string       `json:"updated_at"`
	OpenedAt        string       `json:"opened_at"`
	ClosedAt        string       `json:"closed_at"`
	TeacherFullname string       `json:"teacher_fullname"`
	TeacherID       uint         `json:"teacher_id"`
	Homeworks       HomeworkDtos `json:"homeworks"`
}

type HomeworkPageForm struct {
	Title           string `json:"title" form:"required,alpha_space,max=255"`
	Content         string `json:"content" form:"required,max=255"`
	CourseTitle     string `json:"course_title" form:"required,max=255"`
	OpenedAt        string `json:"opened_at" form:"required,date"`
	ClosedAt        string `json:"closed_at" form:"required,date"`
	TeacherFullname string `json:"teacher_fullname" form:"alpha_space,max=255"`
}

func (hwp HomeworkPage) ToDto() *HomeworkPageDto {
	return &HomeworkPageDto{
		ID:              hwp.ID,
		Title:           hwp.Title,
		Content:         hwp.Content,
		StudentLink:     hwp.StudentLink,
		TeacherLink:     hwp.TeacherLink,
		CourseTitle:     hwp.CourseTitle,
		CreatedAt:       hwp.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       hwp.UpdatedAt.Format(time.RFC3339),
		OpenedAt:        hwp.OpenedAt.Format(time.RFC3339),
		ClosedAt:        hwp.ClosedAt.Format(time.RFC3339),
		TeacherFullname: hwp.TeacherFullname,
		TeacherID:       hwp.TeacherID,
	}
}

func (hwp HomeworkPage) ToNestedDto(hws Homeworks) *HomeworkPageNestedDto {
	return &HomeworkPageNestedDto{
		ID:              hwp.ID,
		Title:           hwp.Title,
		Content:         hwp.Content,
		StudentLink:     hwp.StudentLink,
		TeacherLink:     hwp.TeacherLink,
		CourseTitle:     hwp.CourseTitle,
		CreatedAt:       hwp.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       hwp.UpdatedAt.Format(time.RFC3339),
		OpenedAt:        hwp.OpenedAt.Format(time.RFC3339),
		ClosedAt:        hwp.ClosedAt.Format(time.RFC3339),
		TeacherFullname: hwp.TeacherFullname,
		TeacherID:       hwp.TeacherID,
		Homeworks:       hws.ToDto(),
	}
}

func (hwps HomeworkPages) ToDto() HomeworkPageDtos {
	dtos := make([]*HomeworkPageDto, len(hwps))
	for i, hw := range hwps {
		dtos[i] = hw.ToDto()
	}
	return dtos
}

func (f *HomeworkPageForm) ToModel() (*HomeworkPage, error) {
	openedAt, err := time.Parse(time.RFC3339, f.OpenedAt)
	if err != nil {
		return nil, err
	}

	closedAt, err := time.Parse(time.RFC3339, f.ClosedAt)
	if err != nil {
		return nil, err
	}

	return &HomeworkPage{
		Title:           f.Title,
		Content:         f.Content,
		CourseTitle:     f.CourseTitle,
		OpenedAt:        openedAt,
		ClosedAt:        closedAt,
		TeacherFullname: f.TeacherFullname,
	}, nil
}
