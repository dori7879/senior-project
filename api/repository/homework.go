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

func ReadHomework(db *gorm.DB, id uint) (*model.Homework, error) {
	homework := &model.Homework{}
	if err := db.First(&homework, id).Error; err != nil {
		return nil, err
	}

	return homework, nil
}

func DeleteHomework(db *gorm.DB, id uint) error {
	homework := &model.Homework{}
	if err := db.Where("id = ?", id).Delete(&homework).Error; err != nil {
		return err
	}

	return nil
}

func CreateHomework(db *gorm.DB, homework *model.Homework) (*model.Homework, error) {
	if err := db.Create(homework).Error; err != nil {
		return nil, err
	}

	return homework, nil
}

func UpdateHomework(db *gorm.DB, homework *model.Homework) error {
	if err := db.First(&model.Homework{}, homework.ID).Update(homework).Error; err != nil {
		return err
	}

	return nil
}
