package model

import (
	"time"
)

type HomeworkPages []*HomeworkPage

type HomeworkPage struct {
	ID              uint
	Title           string
	Content         string
	StudentLink     string
	TeacherLink     string
	CourseTitle     string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	OpenedAt        time.Time
	ClosedAt        time.Time
	TeacherFullName string
	TeacherID       uint
	Homeworks       []Homework
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
	TeacherFullName string `json:"teacher_fullname"`
	TeacherID       uint   `json:"teacher_id"`
}

type HomeworkPageForm struct {
	Title           string `json:"title" form:"required,alpha_space,max=255"`
	Content         string `json:"content" form:"required,max=255"`
	CourseTitle     string `json:"course_title" form:"required,max=255"`
	OpenedAt        string `json:"opened_at" form:"required,date"`
	ClosedAt        string `json:"closed_at" form:"required,date"`
	TeacherFullName string `json:"teacher_fullname" form:"required,alpha_space,max=255"`
}

func (hwp HomeworkPage) ToDto() *HomeworkPageDto {
	return &HomeworkPageDto{
		ID:              hwp.ID,
		Title:           hwp.Title,
		Content:         hwp.Content,
		StudentLink:     hwp.StudentLink,
		TeacherLink:     hwp.TeacherLink,
		CourseTitle:     hwp.CourseTitle,
		CreatedAt:       hwp.CreatedAt.Format("2006-01-02"),
		UpdatedAt:       hwp.UpdatedAt.Format("2006-01-02"),
		OpenedAt:        hwp.OpenedAt.Format("2006-01-02"),
		ClosedAt:        hwp.ClosedAt.Format("2006-01-02"),
		TeacherFullName: hwp.TeacherFullName,
		TeacherID:       hwp.TeacherID,
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
	openedAt, err := time.Parse("2008-01-23", f.OpenedAt)
	if err != nil {
		return nil, err
	}

	closedAt, err := time.Parse("2008-01-23", f.ClosedAt)
	if err != nil {
		return nil, err
	}

	return &HomeworkPage{
		Title:           f.Title,
		Content:         f.Content,
		CourseTitle:     f.CourseTitle,
		OpenedAt:        openedAt,
		ClosedAt:        closedAt,
		TeacherFullName: f.TeacherFullName,
	}, nil
}