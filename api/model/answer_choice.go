package model

type AnswerChoices []*AnswerChoice

type AnswerChoice struct {
	ID            uint `gorm:"primaryKey"`
	Content       string
	CorrectAnswer bool
	Fixed         bool
	QuestionID    uint
}

type AnswerChoiceDtos []*AnswerChoiceDto

type AnswerChoiceDto struct {
	ID            uint   `json:"id"`
	Content       string `json:"content"`
	CorrectAnswer bool   `json:"correct_answer"`
	Fixed         bool   `json:"fixed"`
	QuestionID    uint   `json:"question_id"`
}

type AnswerChoiceForm struct {
	Content       string `json:"content" form:"max=255"`
	CorrectAnswer bool   `json:"correct_answer" form:""`
	Fixed         bool   `json:"fixed" form:""`
	QuestionID    uint   `json:"question_id" form:""`
}

func (ac AnswerChoice) ToDto() *AnswerChoiceDto {
	return &AnswerChoiceDto{
		ID:            ac.ID,
		Content:       ac.Content,
		CorrectAnswer: ac.CorrectAnswer,
		Fixed:         ac.Fixed,
		QuestionID:    ac.QuestionID,
	}
}

func (acs AnswerChoices) ToDto() AnswerChoiceDtos {
	dtos := make([]*AnswerChoiceDto, len(acs))
	for i, ac := range acs {
		dtos[i] = ac.ToDto()
	}
	return dtos
}

func (f *AnswerChoiceForm) ToModel() (*AnswerChoice, error) {
	if f.Content != "" {
		return &AnswerChoice{
			Content:       f.Content,
			CorrectAnswer: f.CorrectAnswer,
			QuestionID:    f.QuestionID,
		}, nil
	}

	return &AnswerChoice{
		CorrectAnswer: f.CorrectAnswer,
		QuestionID:    f.QuestionID,
	}, nil
}
