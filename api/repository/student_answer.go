package repository

import (
	"api/model"

	"github.com/jinzhu/gorm"
)

func ListStudentAnswersByQuizSubmission(db *gorm.DB, qsID uint) ([]model.StudentAnswer, error) {
	answers := make([]model.StudentAnswer, 0)
	if err := db.Where("quiz_submission_id = ?", qsID).Find(&answers).Error; err != nil {
		return nil, err
	}

	return answers, nil
}

func ListStudentAnswersByOpenQuestion(db *gorm.DB, oqID uint) (model.StudentAnswers, error) {
	answers := make([]*model.StudentAnswer, 0)
	if err := db.Where("open_question_id = ?", oqID).Find(&answers).Error; err != nil {
		return nil, err
	}

	return answers, nil
}

func ListStudentAnswersByTrueFalseQuestion(db *gorm.DB, tfqID uint) (model.StudentAnswers, error) {
	answers := make([]*model.StudentAnswer, 0)
	if err := db.Where("true_false_question_id = ?", tfqID).Find(&answers).Error; err != nil {
		return nil, err
	}

	return answers, nil
}

func ListStudentAnswersByMultipleChoiceQuestion(db *gorm.DB, mcqID uint) (model.StudentAnswers, error) {
	answers := make([]*model.StudentAnswer, 0)
	if err := db.Where("multple_choice_question_id = ?", mcqID).Find(&answers).Error; err != nil {
		return nil, err
	}

	return answers, nil
}

func DeleteStudentAnswer(db *gorm.DB, id uint) error {
	answer := &model.StudentAnswer{}
	if err := db.Where("id = ?", id).Delete(&answer).Error; err != nil {
		return err
	}

	return nil
}

func CreateStudentAnswer(db *gorm.DB, answer *model.StudentAnswer) (*model.StudentAnswer, error) {
	if answer.Type == "multiple" {
		answer.IsCorrect = false
		acs, err := ListRelatedAnswerChoices(db, answer.MultipleChoiceQuestionID)
		if err != nil {
			return nil, err
		}

		contains := false
		for _, ac := range acs {
			if ac.ID == answer.MultipleChoiceAnswer {
				contains = true
			}
		}

		if contains {
			answer.IsCorrect = true
		}
	} else if answer.Type == "truefalse" {
		tfq, err := ReadTrueFalseQuestion(db, answer.TrueFalseQuestionID)
		if err != nil {
			return nil, err
		}

		if answer.TrueFalseAnswer == tfq.Answer {
			answer.IsCorrect = true
		} else {
			answer.IsCorrect = false
		}
	}

	if err := db.Create(answer).Error; err != nil {
		return nil, err
	}

	return answer, nil
}

func UpdateStudentAnswer(db *gorm.DB, answer *model.StudentAnswer) error {
	if err := db.Model(&model.StudentAnswer{}).Update(answer).Error; err != nil {
		return err
	}

	return nil
}
