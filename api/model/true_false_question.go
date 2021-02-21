package model

import (
	"time"
)

type TrueFalseQuestions []*TrueFalseQuestion

type TrueFalseQuestion struct {
	ID             uint `gorm:"primaryKey"`
	Content        string
	Answer         bool
	Fixed          bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	QuizID         uint
	StudentAnswers []StudentAnswer `gorm:"foreignKey:TrueFalseQuestionID;references:ID"`
}

type TrueFalseQuestionDtos []*TrueFalseQuestionDto

type TrueFalseQuestionDto struct {
	ID        uint   `json:"id"`
	Content   string `json:"content"`
	Answer    bool   `json:"answer"`
	Fixed     bool   `json:"fixed"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	QuizID    uint   `json:"quiz_id"`
}

type TrueFalseQuestionNestedDto struct {
	ID             uint              `json:"id"`
	Content        string            `json:"content"`
	Answer         bool              `json:"answer"`
	Fixed          bool              `json:"fixed"`
	CreatedAt      string            `json:"created_at"`
	UpdatedAt      string            `json:"updated_at"`
	QuizID         uint              `json:"quiz_id"`
	StudentAnswers StudentAnswerDtos `json:"student_answers"`
}

type TrueFalseQuestionForm struct {
	Content string `json:"content" form:"max=255"`
	Answer  bool   `json:"answer" form:""`
	Fixed   bool   `json:"fixed" form:""`
	QuizID  uint   `json:"quiz_id" form:""`
}

func (tfq TrueFalseQuestion) ToDto() *TrueFalseQuestionDto {
	return &TrueFalseQuestionDto{
		ID:      tfq.ID,
		Content: tfq.Content,
	}
}

func (tfq TrueFalseQuestion) ToNestedDto(sts StudentAnswers) *TrueFalseQuestionNestedDto {
	return &TrueFalseQuestionNestedDto{
		ID:             tfq.ID,
		Content:        tfq.Content,
		Answer:         tfq.Answer,
		Fixed:          tfq.Fixed,
		CreatedAt:      tfq.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      tfq.UpdatedAt.Format(time.RFC3339),
		QuizID:         tfq.QuizID,
		StudentAnswers: sts.ToDto(),
	}
}

func (tfqs TrueFalseQuestions) ToDto() TrueFalseQuestionDtos {
	dtos := make([]*TrueFalseQuestionDto, len(tfqs))
	for i, tfq := range tfqs {
		dtos[i] = tfq.ToDto()
	}
	return dtos
}

func (f *TrueFalseQuestionForm) ToModel() (*TrueFalseQuestion, error) {
	if f.Content != "" {
		return &TrueFalseQuestion{
			CreatedAt: time.Now(),
			Content:   f.Content,
			Answer:    f.Answer,
			Fixed:     f.Fixed,
			QuizID:    f.QuizID,
		}, nil
	}

	return &TrueFalseQuestion{
		CreatedAt: time.Now(),
		Answer:    f.Answer,
		Fixed:     f.Fixed,
		QuizID:    f.QuizID,
	}, nil
}
