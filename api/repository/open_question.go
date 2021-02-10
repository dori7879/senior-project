package repository

import (
	"api/model"

	"github.com/jinzhu/gorm"
)

func ListRelatedOpenQuestions(db *gorm.DB, qID uint) (model.OpenQuestions, error) {
	oQuestions := make([]*model.OpenQuestion, 0)
	if err := db.Where("quiz_id = ?", qID).Find(&oQuestions).Error; err != nil {
		return nil, err
	}

	return oQuestions, nil
}

func DeleteOpenQuestionByTeacher(db *gorm.DB, id uint, email string) error {
	teacher := &model.User{}
	if err := db.Where("email = ?", email).First(&teacher).Error; err != nil {
		return err
	}

	quiz := &model.Quiz{}
	if err := db.Where("teacher_id = ?", teacher.ID).First(&quiz).Error; err != nil {
		return err
	}

	oQuestion := &model.OpenQuestion{}
	if err := db.Where("quiz_id = ? AND id = ?", quiz.ID, id).Delete(&oQuestion).Error; err != nil {
		return err
	}

	return nil
}

func CreateOpenQuestion(db *gorm.DB, oQuestion *model.OpenQuestion) (*model.OpenQuestion, error) {
	if err := db.Create(oQuestion).Error; err != nil {
		return nil, err
	}

	return oQuestion, nil
}

func UpdateOpenQuestionByTeacher(db *gorm.DB, oQuestion *model.OpenQuestion, email string) error {
	teacher := &model.User{}
	if err := db.Where("email = ?", email).First(&teacher).Error; err != nil {
		return err
	}

	quiz := &model.Quiz{}
	if err := db.Where("teacher_id = ? AND id = ?", teacher.ID, oQuestion.QuizID).First(&quiz).Error; err != nil {
		return err
	}

	if err := db.Model(&model.OpenQuestion{}).Where("quiz_id = ?", quiz.ID).Update(oQuestion).Error; err != nil {
		return err
	}

	return nil
}
