package model

import (
	"time"
)

type MultipleChoiceQuestions []*MultipleChoiceQuestion

type MultipleChoiceQuestion struct {
	ID             uint `gorm:"primaryKey"`
	Content        string
	Fixed          bool
	SubmittedAt    time.Time
	UpdatedAt      time.Time
	QuizID         uint
	AnswerChoices  []AnswerChoice  `gorm:"foreignKey:QuestionID;references:ID"`
	StudentAnswers []StudentAnswer `gorm:"foreignKey:MultipleChoiceQuestionID;references:ID"`
}

type MultipleChoiceQuestionDtos []*MultipleChoiceQuestionDto

type MultipleChoiceQuestionDto struct {
	ID          uint   `json:"id"`
	Content     string `json:"content"`
	Fixed       bool   `json:"fixed"`
	SubmittedAt string `json:"submitted_at"`
	UpdatedAt   string `json:"updated_at"`
	QuizID      uint   `json:"quiz_id"`
}

type MultipleChoiceQuestionNestedDto struct {
	ID             uint              `json:"id"`
	Content        string            `json:"content"`
	Fixed          bool              `json:"fixed"`
	SubmittedAt    string            `json:"submitted_at"`
	UpdatedAt      string            `json:"updated_at"`
	QuizID         uint              `json:"quiz_id"`
	AnswerChoices  AnswerChoiceDtos  `json:"answer_choices"`
	StudentAnswers StudentAnswerDtos `json:"student_answers"`
}

type MultipleChoiceQuestionForm struct {
	Content string `json:"content" form:"max=255"`
	Fixed   bool   `json:"fixed" form:""`
	QuizID  uint   `json:"quiz_id" form:""`
}

func (mcq MultipleChoiceQuestion) ToDto() *MultipleChoiceQuestionDto {
	return &MultipleChoiceQuestionDto{
		ID:          mcq.ID,
		Content:     mcq.Content,
		Fixed:       mcq.Fixed,
		SubmittedAt: mcq.SubmittedAt.Format(time.RFC3339),
		UpdatedAt:   mcq.UpdatedAt.Format(time.RFC3339),
		QuizID:      mcq.QuizID,
	}
}

func (mcq MultipleChoiceQuestion) ToNestedDto(acs AnswerChoices, sts StudentAnswers) *MultipleChoiceQuestionNestedDto {
	return &MultipleChoiceQuestionNestedDto{
		ID:             mcq.ID,
		Content:        mcq.Content,
		Fixed:          mcq.Fixed,
		SubmittedAt:    mcq.SubmittedAt.Format(time.RFC3339),
		UpdatedAt:      mcq.UpdatedAt.Format(time.RFC3339),
		QuizID:         mcq.QuizID,
		AnswerChoices:  acs.ToDto(),
		StudentAnswers: sts.ToDto(),
	}
}

func (mcqs MultipleChoiceQuestions) ToDto() MultipleChoiceQuestionDtos {
	dtos := make([]*MultipleChoiceQuestionDto, len(mcqs))
	for i, mcq := range mcqs {
		dtos[i] = mcq.ToDto()
	}
	return dtos
}

func (f *MultipleChoiceQuestionForm) ToModel() (*MultipleChoiceQuestion, error) {
	if f.Content != "" {
		return &MultipleChoiceQuestion{
			SubmittedAt: time.Now(),
			Content:     f.Content,
			QuizID:      f.QuizID,
		}, nil
	}

	return &MultipleChoiceQuestion{
		SubmittedAt: time.Now(),
		QuizID:      f.QuizID,
	}, nil
}
