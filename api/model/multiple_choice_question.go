package model

import (
	"time"
)

type MultipleChoiceQuestions []*MultipleChoiceQuestion

type MultipleChoiceQuestion struct {
	ID            uint `gorm:"primaryKey"`
	Content       string
	Fixed         bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
	QuizID        uint
	AnswerChoices []AnswerChoice `gorm:"foreignKey:QuestionID;references:ID"`
}

type MultipleChoiceQuestionDtos []*MultipleChoiceQuestionDto

type MultipleChoiceQuestionDto struct {
	ID        uint   `json:"id"`
	Content   string `json:"content"`
	Fixed     bool   `json:"fixed"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	QuizID    uint   `json:"quiz_id"`
}

type MultipleChoiceQuestionNestedDtos []*MultipleChoiceQuestionNestedDto

type MultipleChoiceQuestionNestedDto struct {
	ID            uint             `json:"id"`
	Content       string           `json:"content"`
	Fixed         bool             `json:"fixed"`
	CreatedAt     string           `json:"created_at"`
	UpdatedAt     string           `json:"updated_at"`
	QuizID        uint             `json:"quiz_id"`
	AnswerChoices AnswerChoiceDtos `json:"answer_choices"`
	// StudentAnswers StudentAnswerDtos `json:"student_answers"`
}

type MultipleChoiceQuestionForm struct {
	Content       string              `json:"content" form:"max=255"`
	Fixed         bool                `json:"fixed" form:""`
	QuizID        uint                `json:"quiz_id" form:""`
	AnswerChoices []*AnswerChoiceForm `json:"answer_choices"`
}

func (mcq MultipleChoiceQuestion) ToDto() *MultipleChoiceQuestionDto {
	return &MultipleChoiceQuestionDto{
		ID:      mcq.ID,
		Content: mcq.Content,
	}
}

func (mcq MultipleChoiceQuestion) ToNestedDto(acs AnswerChoices) *MultipleChoiceQuestionNestedDto {
	return &MultipleChoiceQuestionNestedDto{
		ID:            mcq.ID,
		Content:       mcq.Content,
		Fixed:         mcq.Fixed,
		CreatedAt:     mcq.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     mcq.UpdatedAt.Format(time.RFC3339),
		QuizID:        mcq.QuizID,
		AnswerChoices: acs.ToDto(),
		// StudentAnswers: sts.ToDto(),
	}
}

func (mcqs MultipleChoiceQuestions) ToDto() MultipleChoiceQuestionDtos {
	dtos := make([]*MultipleChoiceQuestionDto, len(mcqs))
	for i, mcq := range mcqs {
		dtos[i] = mcq.ToDto()
	}
	return dtos
}

func (mcqs MultipleChoiceQuestions) ToNestedDto() MultipleChoiceQuestionNestedDtos {
	dtos := make([]*MultipleChoiceQuestionNestedDto, len(mcqs))
	for i, mcq := range mcqs {
		answerChoices := make(AnswerChoices, len(mcq.AnswerChoices))
		for i, v := range mcq.AnswerChoices {
			tempV := v
			answerChoices[i] = &tempV
		}
		dtos[i] = mcq.ToNestedDto(answerChoices)
	}
	return dtos
}

func (f *MultipleChoiceQuestionForm) ToModel() (*MultipleChoiceQuestion, error) {
	if f.Content != "" {
		return &MultipleChoiceQuestion{
			CreatedAt: time.Now(),
			Content:   f.Content,
			QuizID:    f.QuizID,
			Fixed:     f.Fixed,
		}, nil
	}

	return &MultipleChoiceQuestion{
		CreatedAt: time.Now(),
		QuizID:    f.QuizID,
		Fixed:     f.Fixed,
	}, nil
}
