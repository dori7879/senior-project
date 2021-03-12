package model

import (
	"time"
)

type Quizzes []*Quiz

type Quiz struct {
	ID                      uint `gorm:"primaryKey"`
	Title                   string
	Content                 string
	StudentLink             string `gorm:"unique"`
	TeacherLink             string `gorm:"unique"`
	CourseTitle             string
	Mode                    string
	CreatedAt               time.Time
	UpdatedAt               time.Time
	OpenedAt                time.Time
	ClosedAt                time.Time
	TeacherFullname         string
	TeacherID               uint                     `gorm:"default:null"`
	StudentGroupID          uint                     `gorm:"default:null"`
	QuizSubmissions         []QuizSubmission         `gorm:"foreignKey:QuizID;references:ID"`
	OpenQuestions           []OpenQuestion           `gorm:"foreignKey:QuizID;references:ID"`
	TrueFalseQuestions      []TrueFalseQuestion      `gorm:"foreignKey:QuizID;references:ID"`
	MultipleChoiceQuestions []MultipleChoiceQuestion `gorm:"foreignKey:QuizID;references:ID"`
}

type QuizDtos []*QuizDto

type QuizDto struct {
	ID              uint   `json:"id"`
	Title           string `json:"title"`
	Content         string `json:"content"`
	StudentLink     string `json:"student_link"`
	TeacherLink     string `json:"teacher_link"`
	CourseTitle     string `json:"course_title"`
	Mode            string `json:"mode"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	OpenedAt        string `json:"opened_at"`
	ClosedAt        string `json:"closed_at"`
	TeacherFullname string `json:"teacher_fullname"`
	TeacherID       uint   `json:"teacher_id"`
}

type QuizNestedDto struct {
	ID                      uint                             `json:"id"`
	Title                   string                           `json:"title"`
	Content                 string                           `json:"content"`
	StudentLink             string                           `json:"student_link"`
	TeacherLink             string                           `json:"teacher_link"`
	CourseTitle             string                           `json:"course_title"`
	Mode                    string                           `json:"mode"`
	CreatedAt               string                           `json:"created_at"`
	UpdatedAt               string                           `json:"updated_at"`
	OpenedAt                string                           `json:"opened_at"`
	ClosedAt                string                           `json:"closed_at"`
	TeacherFullname         string                           `json:"teacher_fullname"`
	TeacherID               uint                             `json:"teacher_id"`
	QuizSubmissions         QuizSubmissionNestedDtos         `json:"quiz_submissions"`
	OpenQuestions           OpenQuestionDtos                 `json:"open_questions"`
	TrueFalseQuestions      TrueFalseQuestionDtos            `json:"true_false_questions"`
	MultipleChoiceQuestions MultipleChoiceQuestionNestedDtos `json:"multiple_choice_questions"`
}

type QuizForm struct {
	Title                   string                        `json:"title" form:"required,alpha_space,max=255"`
	Content                 string                        `json:"content" form:"required,max=255"`
	CourseTitle             string                        `json:"course_title" form:"required,max=255"`
	Mode                    string                        `json:"mode"`
	OpenedAt                string                        `json:"opened_at" form:"required,date"`
	ClosedAt                string                        `json:"closed_at" form:"required,date"`
	TeacherFullname         string                        `json:"teacher_fullname" form:"alpha_space,max=255"`
	OpenQuestions           []*OpenQuestionForm           `json:"open_questions"`
	TrueFalseQuestions      []*TrueFalseQuestionForm      `json:"true_false_questions"`
	MultipleChoiceQuestions []*MultipleChoiceQuestionForm `json:"multiple_choice_questions"`
}

func (q Quiz) ToDto() *QuizDto {
	return &QuizDto{
		ID:              q.ID,
		Title:           q.Title,
		Content:         q.Content,
		StudentLink:     q.StudentLink,
		TeacherLink:     q.TeacherLink,
		CourseTitle:     q.CourseTitle,
		Mode:            q.Mode,
		CreatedAt:       q.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       q.UpdatedAt.Format(time.RFC3339),
		OpenedAt:        q.OpenedAt.Format(time.RFC3339),
		ClosedAt:        q.ClosedAt.Format(time.RFC3339),
		TeacherFullname: q.TeacherFullname,
		TeacherID:       q.TeacherID,
	}
}

func (q Quiz) ToNestedDto(qss QuizSubmissions, oqs OpenQuestions, tfqs TrueFalseQuestions, mcqs MultipleChoiceQuestions) *QuizNestedDto {
	return &QuizNestedDto{
		ID:                      q.ID,
		Title:                   q.Title,
		Content:                 q.Content,
		StudentLink:             q.StudentLink,
		TeacherLink:             q.TeacherLink,
		CourseTitle:             q.CourseTitle,
		Mode:                    q.Mode,
		CreatedAt:               q.CreatedAt.Format(time.RFC3339),
		UpdatedAt:               q.UpdatedAt.Format(time.RFC3339),
		OpenedAt:                q.OpenedAt.Format(time.RFC3339),
		ClosedAt:                q.ClosedAt.Format(time.RFC3339),
		TeacherFullname:         q.TeacherFullname,
		TeacherID:               q.TeacherID,
		QuizSubmissions:         qss.ToNestedDto(),
		OpenQuestions:           oqs.ToDto(),
		TrueFalseQuestions:      tfqs.ToDto(),
		MultipleChoiceQuestions: mcqs.ToNestedDto(),
	}
}

func (qs Quizzes) ToDto() QuizDtos {
	dtos := make([]*QuizDto, len(qs))
	for i, q := range qs {
		dtos[i] = q.ToDto()
	}
	return dtos
}

func (f *QuizForm) ToModel() (*Quiz, error) {
	openedAt, err := time.Parse(time.RFC3339, f.OpenedAt)
	if err != nil {
		return nil, err
	}

	closedAt, err := time.Parse(time.RFC3339, f.ClosedAt)
	if err != nil {
		return nil, err
	}

	return &Quiz{
		Title:           f.Title,
		Content:         f.Content,
		CourseTitle:     f.CourseTitle,
		Mode:            f.Mode,
		OpenedAt:        openedAt,
		ClosedAt:        closedAt,
		TeacherFullname: f.TeacherFullname,
	}, nil
}
