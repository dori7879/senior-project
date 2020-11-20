package repository

import (
	"api/model"

	"github.com/jinzhu/gorm"
)

func GetUser(db *gorm.DB, id uint) (*model.User, error) {
	user := &model.User{}
	if err := db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByEmail(db *gorm.DB, email string) (*model.User, error) {
	user := &model.User{}
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func CreateUser(db *gorm.DB, user *model.User) (*model.User, error) {
	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func UpdateUser(db *gorm.DB, user *model.User) error {
	if err := db.First(&model.User{}, user.ID).Update(user).Error; err != nil {
		return err
	}

	return nil
}
