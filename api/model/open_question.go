package model

import (
	"time"
)

type OpenQuestions []*OpenQuestion

type OpenQuestion struct {
	ID             uint `gorm:"primaryKey"`
	Content        string
	Answer         string
	Fixed          bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	QuizID         uint
	StudentAnswers []StudentAnswer `gorm:"foreignKey:OpenQuestionID;references:ID"`
}

type OpenQuestionDtos []*OpenQuestionDto

type OpenQuestionDto struct {
	ID        uint   `json:"id"`
	Content   string `json:"content"`
	Answer    string `json:"answer"`
	Fixed     bool   `json:"fixed"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	QuizID    uint   `json:"quiz_id"`
}

type OpenQuestionNestedDto struct {
	ID             uint              `json:"id"`
	Content        string            `json:"content"`
	Answer         string            `json:"answer"`
	Fixed          bool              `json:"fixed"`
	CreatedAt      string            `json:"created_at"`
	UpdatedAt      string            `json:"updated_at"`
	QuizID         uint              `json:"quiz_id"`
	StudentAnswers StudentAnswerDtos `json:"student_answers"`
}

type OpenQuestionForm struct {
	Content string `json:"content" form:"max=255"`
	Answer  string `json:"answer" form:"max=255"`
	Fixed   bool   `json:"fixed" form:""`
	QuizID  uint   `json:"quiz_id" form:""`
}

func (oq OpenQuestion) ToDto() *OpenQuestionDto {
	return &OpenQuestionDto{
		ID:      oq.ID,
		Content: oq.Content,
	}
}

func (oq OpenQuestion) ToNestedDto(sts StudentAnswers) *OpenQuestionNestedDto {
	return &OpenQuestionNestedDto{
		ID:             oq.ID,
		Content:        oq.Content,
		Answer:         oq.Answer,
		Fixed:          oq.Fixed,
		CreatedAt:      oq.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      oq.UpdatedAt.Format(time.RFC3339),
		QuizID:         oq.QuizID,
		StudentAnswers: sts.ToDto(),
	}
}

func (oqs OpenQuestions) ToDto() OpenQuestionDtos {
	dtos := make([]*OpenQuestionDto, len(oqs))
	for i, oq := range oqs {
		dtos[i] = oq.ToDto()
	}
	return dtos
}

func (f *OpenQuestionForm) ToModel() (*OpenQuestion, error) {
	if f.Content != "" || f.Answer != "" {
		return &OpenQuestion{
			CreatedAt: time.Now(),
			Content:   f.Content,
			Answer:    f.Answer,
			Fixed:     f.Fixed,
			QuizID:    f.QuizID,
		}, nil
	}

	return &OpenQuestion{
		CreatedAt: time.Now(),
		Fixed:     f.Fixed,
		QuizID:    f.QuizID,
	}, nil
}
