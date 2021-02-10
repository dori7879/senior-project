package repository

import (
	"api/model"

	"github.com/jinzhu/gorm"
)

func ListRelatedAnswerChoices(db *gorm.DB, mcqID uint) (model.AnswerChoices, error) {
	answerChoices := make([]*model.AnswerChoice, 0)
	if err := db.Where("question_id = ?", mcqID).Find(&answerChoices).Error; err != nil {
		return nil, err
	}

	return answerChoices, nil
}

func DeleteAnswerChoice(db *gorm.DB, id uint) error {
	answerChoice := &model.AnswerChoice{}
	if err := db.Where("id = ?", id).Delete(&answerChoice).Error; err != nil {
		return err
	}

	return nil
}

func CreateAnswerChoice(db *gorm.DB, answerChoice *model.AnswerChoice) (*model.AnswerChoice, error) {
	if err := db.Create(answerChoice).Error; err != nil {
		return nil, err
	}

	return answerChoice, nil
}

func UpdateAnswerChoice(db *gorm.DB, answerChoice *model.AnswerChoice) error {
	if err := db.Model(&model.AnswerChoice{}).Update(answerChoice).Error; err != nil {
		return err
	}

	return nil
}
