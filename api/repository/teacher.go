package repository

import (
	"api/model"

	"github.com/jinzhu/gorm"
)

func GetTeacher(db *gorm.DB, id uint) (*model.Teacher, error) {
	teacher := &model.Teacher{}
	if err := db.Where("teacher_id = ?", id).First(&teacher).Error; err != nil {
		return nil, err
	}

	return teacher, nil
}

func GetTeacherByEmail(db *gorm.DB, email string) (*model.Teacher, error) {
	teacher := &model.Teacher{}
	if err := db.Where("email = ?", email).First(&teacher).Error; err != nil {
		return nil, err
	}

	return teacher, nil
}
