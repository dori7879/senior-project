package repository

import (
	"api/model"

	"github.com/jinzhu/gorm"
)

func ListRelatedMultipleChoiceQuestions(db *gorm.DB, qID uint) (model.MultipleChoiceQuestions, error) {
	mcQuestions := make([]*model.MultipleChoiceQuestion, 0)
	if err := db.Where("quiz_id = ?", qID).Find(&mcQuestions).Error; err != nil {
		return nil, err
	}

	return mcQuestions, nil
}

func DeleteMultipleChoiceQuestionByTeacher(db *gorm.DB, id uint, email string) error {
	teacher := &model.User{}
	if err := db.Where("email = ?", email).First(&teacher).Error; err != nil {
		return err
	}

	quiz := &model.Quiz{}
	if err := db.Where("teacher_id = ?", teacher.ID).First(&quiz).Error; err != nil {
		return err
	}

	mcQuestion := &model.MultipleChoiceQuestion{}
	if err := db.Where("quiz_id = ? AND id = ?", quiz.ID, id).Delete(&mcQuestion).Error; err != nil {
		return err
	}

	return nil
}

func CreateMultipleChoiceQuestion(db *gorm.DB, mcQuestion *model.MultipleChoiceQuestion) (*model.MultipleChoiceQuestion, error) {
	if err := db.Create(mcQuestion).Error; err != nil {
		return nil, err
	}

	return mcQuestion, nil
}

func UpdateMultipleChoiceQuestionByTeacher(db *gorm.DB, mcQuestion *model.MultipleChoiceQuestion, email string) error {
	teacher := &model.User{}
	if err := db.Where("email = ?", email).First(&teacher).Error; err != nil {
		return err
	}

	quiz := &model.Quiz{}
	if err := db.Where("teacher_id = ? AND id = ?", teacher.ID, mcQuestion.QuizID).First(&quiz).Error; err != nil {
		return err
	}

	if err := db.Model(&model.MultipleChoiceQuestion{}).Where("quiz_id = ?", quiz.ID).Update(mcQuestion).Error; err != nil {
		return err
	}

	return nil
}
