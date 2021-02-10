package repository

import (
	"api/model"

	"github.com/jinzhu/gorm"
)

func ListQuizzesByOwner(db *gorm.DB, email string) (model.Quizzes, error) {
	user := &model.User{}
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	quizzes := make([]*model.Quiz, 0)
	if err := db.Where("teacher_id = ?", user.ID).Find(&quizzes).Error; err != nil {
		return nil, err
	}

	return quizzes, nil
}

func ReadQuiz(db *gorm.DB, id uint) (*model.Quiz, error) {
	quiz := &model.Quiz{}
	if err := db.Where("id = ?", id).First(&quiz).Error; err != nil {
		return nil, err
	}

	return quiz, nil
}

func ReadQuizWithNoOwner(db *gorm.DB, id uint) (*model.Quiz, error) {
	quiz := &model.Quiz{}
	if err := db.Where("teacher_id is NULL AND id = ?", id).First(&quiz).Error; err != nil {
		return nil, err
	}

	return quiz, nil
}

func ReadQuizWithNoOwnerByStudentLink(db *gorm.DB, link string) (*model.Quiz, error) {
	quiz := &model.Quiz{}
	if err := db.Where("teacher_id is NULL AND student_link = ?", link).First(&quiz).Error; err != nil {
		return nil, err
	}

	return quiz, nil
}

func ReadQuizWithNoOwnerByTeacherLink(db *gorm.DB, link string) (*model.Quiz, error) {
	quiz := &model.Quiz{}
	if err := db.Where("teacher_id is NULL AND teacher_link = ?", link).First(&quiz).Error; err != nil {
		return nil, err
	}

	return quiz, nil
}

func ReadQuizByOwner(db *gorm.DB, id uint, email string) (*model.Quiz, error) {
	user := &model.User{}
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	quiz := &model.Quiz{}
	if err := db.Where("teacher_id = ? AND id = ?", user.ID, id).First(&quiz).Error; err != nil {
		return nil, err
	}

	return quiz, nil
}

func ReadQuizWithOwnerByTeacherLink(db *gorm.DB, link, email string) (*model.Quiz, error) {
	user := &model.User{}
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	quiz := &model.Quiz{}
	if err := db.Where("teacher_id = ? AND teacher_link = ?", user.ID, link).First(&quiz).Error; err != nil {
		return nil, err
	}

	return quiz, nil
}

func ReadQuizByOwnerByStudentLink(db *gorm.DB, link string) (*model.Quiz, error) {
	quiz := &model.Quiz{}
	if err := db.Where("teacher_id IS NOT NULL AND student_link = ?", link).First(&quiz).Error; err != nil {
		return nil, err
	}

	return quiz, nil
}

func DeleteQuizByOwner(db *gorm.DB, id uint, email string) error {
	user := &model.User{}
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}

	quiz := &model.Quiz{}
	if err := db.Where("id = ? AND teacher_id = ?", id, user.ID).Delete(&quiz).Error; err != nil {
		return err
	}

	return nil
}

func CreateQuiz(db *gorm.DB, quiz *model.Quiz) (*model.Quiz, error) {
	if err := db.Create(quiz).Error; err != nil {
		return nil, err
	}

	return quiz, nil
}

func UpdateQuiz(db *gorm.DB, quiz *model.Quiz) error {
	if err := db.First(&model.Quiz{}, quiz.ID).Update(quiz).Error; err != nil {
		return err
	}

	return nil
}

func UpdateQuizByOwner(db *gorm.DB, quiz *model.Quiz, email string) error {
	user := &model.User{}
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}

	if err := db.Where("teacher_id = ? AND id = ?", user.ID, quiz.ID).First(&model.Quiz{}).Update(quiz).Error; err != nil {
		return err
	}

	return nil
}

func UpdateQuizWithNoOwner(db *gorm.DB, quiz *model.Quiz) error {
	if err := db.Where("teacher_id is NULL AND id = ?", quiz.ID).First(&model.Quiz{}).Update(quiz).Error; err != nil {
		return err
	}

	return nil
}
