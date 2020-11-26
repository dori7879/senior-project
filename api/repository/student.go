package repository

import (
	"api/model"

	"github.com/jinzhu/gorm"
)

func GetStudent(db *gorm.DB, id uint) (*model.Student, error) {
	student := &model.Student{}
	if err := db.Where("student_id = ?", id).First(&student).Error; err != nil {
		return nil, err
	}

	return student, nil
}

func GetStudentByEmail(db *gorm.DB, email string) (*model.Student, error) {
	student := &model.Student{}
	if err := db.Where("email = ?", email).First(&student).Error; err != nil {
		return nil, err
	}

	return student, nil
}
