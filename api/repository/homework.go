package repository

import (
	"api/model"

	"github.com/jinzhu/gorm"
)

func ListHomeworks(db *gorm.DB) (model.Homeworks, error) {
	homeworks := make([]*model.Homework, 0)
	if err := db.Find(&homeworks).Error; err != nil {
		return nil, err
	}

	return homeworks, nil
}
