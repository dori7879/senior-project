package repository

import (
	"api/model"

	"github.com/jinzhu/gorm"
)

func ListRelatedTrueFalseQuestions(db *gorm.DB, qID uint) (model.TrueFalseQuestions, error) {
	tfQuestions := make([]*model.TrueFalseQuestion, 0)
	if err := db.Where("quiz_id = ?", qID).Find(&tfQuestions).Error; err != nil {
		return nil, err
	}

	return tfQuestions, nil
}

func ReadTrueFalseQuestion(db *gorm.DB, id uint) (*model.TrueFalseQuestion, error) {
	tfq := &model.TrueFalseQuestion{}
	if err := db.Where("id = ?", id).First(&tfq).Error; err != nil {
		return nil, err
	}

	return tfq, nil
}

func DeleteTrueFalseQuestionByTeacher(db *gorm.DB, id uint, email string) error {
	teacher := &model.User{}
	if err := db.Where("email = ?", email).First(&teacher).Error; err != nil {
		return err
	}

	quiz := &model.Quiz{}
	if err := db.Where("teacher_id = ?", teacher.ID).First(&quiz).Error; err != nil {
		return err
	}

	tfQuestion := &model.TrueFalseQuestion{}
	if err := db.Where("quiz_id = ? AND id = ?", quiz.ID, id).Delete(&tfQuestion).Error; err != nil {
		return err
	}

	return nil
}

func CreateTrueFalseQuestion(db *gorm.DB, tfQuestion *model.TrueFalseQuestion) (*model.TrueFalseQuestion, error) {
	if err := db.Create(tfQuestion).Error; err != nil {
		return nil, err
	}

	return tfQuestion, nil
}

func UpdateTrueFalseQuestionByTeacher(db *gorm.DB, tfQuestion *model.TrueFalseQuestion, email string) error {
	teacher := &model.User{}
	if err := db.Where("email = ?", email).First(&teacher).Error; err != nil {
		return err
	}

	quiz := &model.Quiz{}
	if err := db.Where("teacher_id = ? AND id = ?", teacher.ID, tfQuestion.QuizID).First(&quiz).Error; err != nil {
		return err
	}

	if err := db.Model(&model.TrueFalseQuestion{}).Where("quiz_id = ?", quiz.ID).Update(tfQuestion).Error; err != nil {
		return err
	}

	return nil
}
