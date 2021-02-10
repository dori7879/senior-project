package model

import (
	"time"
)

type QuizSubmissions []*QuizSubmission

type QuizSubmission struct {
	ID              uint `gorm:"primaryKey"`
	Grade           string
	Comments        string
	SubmittedAt     time.Time `gorm:"default:null"`
	UpdatedAt       time.Time
	StudentFullname string
	StudentID       uint `gorm:"default:null"`
	QuizID          uint
	StudentAnswers  []StudentAnswer `gorm:"foreignKey:QuizSubmissionID;references:ID"`
}

type QuizSubmissionDtos []*QuizSubmissionDto

type QuizSubmissionDto struct {
	ID              uint   `json:"id"`
	Grade           string `json:"grade"`
	Comments        string `json:"comments"`
	SubmittedAt     string `json:"submitted_at"`
	UpdatedAt       string `json:"updated_at"`
	StudentFullname string `json:"student_fullname"`
	StudentID       uint   `json:"student_id"`
	QuizID          uint   `json:"quiz_id"`
}

type QuizSubmissionNestedDto struct {
	ID              uint              `json:"id"`
	Grade           string            `json:"grade"`
	Comments        string            `json:"comments"`
	SubmittedAt     string            `json:"submitted_at"`
	UpdatedAt       string            `json:"updated_at"`
	StudentFullname string            `json:"student_fullname"`
	StudentID       uint              `json:"student_id"`
	QuizID          uint              `json:"quiz_id"`
	StudentAnswers  StudentAnswerDtos `json:"student_answers"`
}

type QuizSubmissionForm struct {
	Grade           string `json:"grade" form:"max=255"`
	Comments        string `json:"comments" form:"max=255"`
	StudentFullname string `json:"student_fullname" form:"alpha_space,max=255"`
	QuizID          uint   `json:"quiz_id" form:""`
}

func (qs QuizSubmission) ToDto() *QuizSubmissionDto {
	return &QuizSubmissionDto{
		ID:              qs.ID,
		Grade:           qs.Grade,
		Comments:        qs.Comments,
		SubmittedAt:     qs.SubmittedAt.Format(time.RFC3339),
		UpdatedAt:       qs.UpdatedAt.Format(time.RFC3339),
		StudentFullname: qs.StudentFullname,
		StudentID:       qs.StudentID,
		QuizID:          qs.QuizID,
	}
}

func (qs QuizSubmission) ToNestedDto(sts StudentAnswers) *QuizSubmissionNestedDto {
	return &QuizSubmissionNestedDto{
		ID:              qs.ID,
		Grade:           qs.Grade,
		Comments:        qs.Comments,
		SubmittedAt:     qs.SubmittedAt.Format(time.RFC3339),
		UpdatedAt:       qs.UpdatedAt.Format(time.RFC3339),
		StudentFullname: qs.StudentFullname,
		StudentID:       qs.StudentID,
		QuizID:          qs.QuizID,
		StudentAnswers:  sts.ToDto(),
	}
}

func (qss QuizSubmissions) ToDto() QuizSubmissionDtos {
	dtos := make([]*QuizSubmissionDto, len(qss))
	for i, qs := range qss {
		dtos[i] = qs.ToDto()
	}
	return dtos
}

func (f *QuizSubmissionForm) ToModel() (*QuizSubmission, error) {
	if f.Grade != "" || f.Comments != "" {
		return &QuizSubmission{
			Grade:           f.Grade,
			Comments:        f.Comments,
			StudentFullname: f.StudentFullname,
			QuizID:          f.QuizID,
		}, nil
	}

	return &QuizSubmission{
		SubmittedAt:     time.Now(),
		StudentFullname: f.StudentFullname,
		QuizID:          f.QuizID,
	}, nil
}
