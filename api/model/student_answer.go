package model

type StudentAnswers []*StudentAnswer

type StudentAnswer struct {
	ID                       uint `gorm:"primaryKey"`
	Type                     string
	OpenAnswer               string
	TrueFalseAnswer          bool
	MultipleChoiceAnswer     uint
	Comments                 string
	QuizSubmissionID         uint
	OpenQuestionID           uint `gorm:"default:null"`
	TrueFalseQuestionID      uint `gorm:"default:null"`
	MultipleChoiceQuestionID uint `gorm:"default:null"`
}

type StudentAnswerDtos []*StudentAnswerDto

type StudentAnswerDto struct {
	ID                       uint   `json:"id"`
	Type                     string `json:"type"`
	OpenAnswer               string `json:"open_answer"`
	TrueFalseAnswer          bool   `json:"true_false_answer"`
	MultipleChoiceAnswer     uint   `json:"multiple_choice_answer"`
	Comments                 string `json:"comments"`
	QuizSubmissionID         uint   `json:"quiz_submission_id"`
	OpenQuestionID           uint   `json:"open_question_id"`
	TrueFalseQuestionID      uint   `json:"true_false_question_id"`
	MultipleChoiceQuestionID uint   `json:"multiple_choice_question_id"`
}

type StudentAnswerForm struct {
	Type                     string `json:"type" form:"max=255"`
	OpenAnswer               string `json:"open_answer" form:""`
	TrueFalseAnswer          bool   `json:"true_false_answer" form:""`
	MultipleChoiceAnswer     uint   `json:"multiple_choice_answer" form:""`
	Comments                 string `json:"comments" form:"max=255"`
	QuizSubmissionID         uint   `json:"quiz_submission_id" form:""`
	OpenQuestionID           uint   `json:"open_question_id" form:""`
	TrueFalseQuestionID      uint   `json:"true_false_question_id" form:""`
	MultipleChoiceQuestionID uint   `json:"multiple_choice_question_id" form:""`
}

func (sa StudentAnswer) ToDto() *StudentAnswerDto {
	return &StudentAnswerDto{
		ID:                       sa.ID,
		Type:                     sa.Type,
		OpenAnswer:               sa.OpenAnswer,
		TrueFalseAnswer:          sa.TrueFalseAnswer,
		MultipleChoiceAnswer:     sa.MultipleChoiceAnswer,
		Comments:                 sa.Comments,
		QuizSubmissionID:         sa.QuizSubmissionID,
		OpenQuestionID:           sa.OpenQuestionID,
		TrueFalseQuestionID:      sa.TrueFalseQuestionID,
		MultipleChoiceQuestionID: sa.MultipleChoiceQuestionID,
	}
}

func (sas StudentAnswers) ToDto() StudentAnswerDtos {
	dtos := make([]*StudentAnswerDto, len(sas))
	for i, sa := range sas {
		dtos[i] = sa.ToDto()
	}
	return dtos
}

func (f *StudentAnswerForm) ToModel() (*StudentAnswer, error) {
	return &StudentAnswer{
		Type:                     f.Type,
		OpenAnswer:               f.OpenAnswer,
		TrueFalseAnswer:          f.TrueFalseAnswer,
		MultipleChoiceAnswer:     f.MultipleChoiceAnswer,
		Comments:                 f.Comments,
		QuizSubmissionID:         f.QuizSubmissionID,
		OpenQuestionID:           f.OpenQuestionID,
		TrueFalseQuestionID:      f.TrueFalseQuestionID,
		MultipleChoiceQuestionID: f.MultipleChoiceQuestionID,
	}, nil
}
