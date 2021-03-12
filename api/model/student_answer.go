package model

import "fmt"

type StudentAnswers []*StudentAnswer

type StudentAnswer struct {
	ID                       uint `gorm:"primaryKey"`
	Type                     string
	OpenAnswer               string
	TrueFalseAnswer          bool
	MultipleChoiceAnswer     uint
	Comments                 string
	IsCorrect                bool
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
	IsCorrect                bool   `json:"is_correct"`
	QuizSubmissionID         uint   `json:"quiz_submission_id"`
	OpenQuestionID           uint   `json:"open_question_id"`
	TrueFalseQuestionID      uint   `json:"true_false_question_id"`
	MultipleChoiceQuestionID uint   `json:"multiple_choice_question_id"`
}

type StudentAnswerForm struct {
	Type                     string `json:"type" form:"max=255"`
	OpenAnswer               string `json:"open_answer" form:""`
	TrueFalseAnswer          bool   `json:"true_false_answer" form:""`
	MultipleChoiceAnswer     []uint `json:"multiple_choice_answer" form:""`
	Comments                 string `json:"comments" form:"max=255"`
	IsCorrect                bool   `json:"is_correct"`
	QuizSubmissionID         uint   `json:"quiz_submission_id" form:""`
	OpenQuestionID           uint   `json:"open_question_id" form:""`
	TrueFalseQuestionID      uint   `json:"true_false_question_id" form:""`
	MultipleChoiceQuestionID uint   `json:"multiple_choice_question_id" form:""`
}

type StudentAnswerUpdateForm struct {
	ID                       uint   `json:"id" form:""`
	Type                     string `json:"type" form:"max=255"`
	OpenAnswer               string `json:"open_answer" form:""`
	TrueFalseAnswer          bool   `json:"true_false_answer" form:""`
	MultipleChoiceAnswer     uint   `json:"multiple_choice_answer" form:""`
	Comments                 string `json:"comments" form:"max=255"`
	IsCorrect                bool   `json:"is_correct"`
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
		IsCorrect:                sa.IsCorrect,
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

func (f *StudentAnswerForm) ToModel() ([]*StudentAnswer, error) {
	var result []*StudentAnswer

	if f.Type != "multiple" {
		var mcAnswer uint
		result = append(result, &StudentAnswer{
			Type:                     f.Type,
			OpenAnswer:               f.OpenAnswer,
			TrueFalseAnswer:          f.TrueFalseAnswer,
			MultipleChoiceAnswer:     mcAnswer,
			Comments:                 f.Comments,
			IsCorrect:                f.IsCorrect,
			QuizSubmissionID:         f.QuizSubmissionID,
			OpenQuestionID:           f.OpenQuestionID,
			TrueFalseQuestionID:      f.TrueFalseQuestionID,
			MultipleChoiceQuestionID: f.MultipleChoiceQuestionID,
		})
		return result, nil
	}

	for _, v := range f.MultipleChoiceAnswer {
		mcAnswer := v
		result = append(result, &StudentAnswer{
			Type:                     f.Type,
			OpenAnswer:               f.OpenAnswer,
			TrueFalseAnswer:          f.TrueFalseAnswer,
			MultipleChoiceAnswer:     mcAnswer,
			Comments:                 f.Comments,
			IsCorrect:                f.IsCorrect,
			QuizSubmissionID:         f.QuizSubmissionID,
			OpenQuestionID:           f.OpenQuestionID,
			TrueFalseQuestionID:      f.TrueFalseQuestionID,
			MultipleChoiceQuestionID: f.MultipleChoiceQuestionID,
		})
	}

	return result, nil
}

func (f *StudentAnswerUpdateForm) ToUpdateModel() (*StudentAnswer, error) {
	fmt.Printf("%T %d\n", f.ID, f.ID)
	fmt.Println(f.ID)
	result := &StudentAnswer{
		ID:                       f.ID,
		Type:                     f.Type,
		OpenAnswer:               f.OpenAnswer,
		TrueFalseAnswer:          f.TrueFalseAnswer,
		MultipleChoiceAnswer:     f.MultipleChoiceAnswer,
		Comments:                 f.Comments,
		IsCorrect:                f.IsCorrect,
		QuizSubmissionID:         f.QuizSubmissionID,
		OpenQuestionID:           f.OpenQuestionID,
		TrueFalseQuestionID:      f.TrueFalseQuestionID,
		MultipleChoiceQuestionID: f.MultipleChoiceQuestionID,
	}

	return result, nil
}
