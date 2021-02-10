package repository

import (
	"api/model"

	"github.com/jinzhu/gorm"
)

func ListQuizSubmissionsByOwnerEmail(db *gorm.DB, email string) (model.QuizSubmissions, error) {
	owner := &model.Student{}
	if err := db.Where("email = ?", email).First(&owner).Error; err != nil {
		return nil, err
	}

	qSubmissions := make([]*model.QuizSubmission, 0)
	if err := db.Where("student_id = ?", owner.ID).Find(&qSubmissions).Error; err != nil {
		return nil, err
	}

	return qSubmissions, nil
}

func ListRelatedQuizSubmissions(db *gorm.DB, qID uint) (model.QuizSubmissions, error) {
	qSubmissions := make([]*model.QuizSubmission, 0)
	if err := db.Where("quiz_id = ?", qID).Find(&qSubmissions).Error; err != nil {
		return nil, err
	}

	return qSubmissions, nil
}

func ReadQuizSubmission(db *gorm.DB, id uint) (*model.QuizSubmission, error) {
	qSubmission := &model.QuizSubmission{}
	if err := db.Where("id = ?", id).First(&qSubmission).Error; err != nil {
		return nil, err
	}

	return qSubmission, nil
}

func ReadQuizSubmissionByIDandOwner(db *gorm.DB, id uint, email string) (*model.QuizSubmission, error) {
	owner := &model.User{}
	if err := db.Where("email = ?", email).First(&owner).Error; err != nil {
		return nil, err
	}

	qSubmission := &model.QuizSubmission{}
	if err := db.Where("student_id = ? AND id = ?", owner.ID, id).First(&qSubmission).Error; err != nil {
		return nil, err
	}

	return qSubmission, nil
}

func ReadQuizSubmissionByIDandTeacher(db *gorm.DB, id uint, email string) (*model.QuizSubmission, error) {
	teacher := &model.User{}
	if err := db.Where("email = ?", email).First(&teacher).Error; err != nil {
		return nil, err
	}

	quiz := &model.Quiz{}
	if err := db.Where("teacher_id = ?", teacher.ID).First(&quiz).Error; err != nil {
		return nil, err
	}

	qSubmission := &model.QuizSubmission{}
	if err := db.Where("id = ? AND quiz_id = ?", id, quiz.ID).First(&qSubmission).Error; err != nil {
		return nil, err
	}

	return qSubmission, nil
}

func DeleteQuizSubmission(db *gorm.DB, id uint) error {
	qSubmission := &model.QuizSubmission{}
	if err := db.Where("id = ?", id).Delete(&qSubmission).Error; err != nil {
		return err
	}

	return nil
}

func DeleteQuizSubmissionByTeacher(db *gorm.DB, id uint, email string) error {
	teacher := &model.User{}
	if err := db.Where("email = ?", email).First(&teacher).Error; err != nil {
		return err
	}

	quiz := &model.Quiz{}
	if err := db.Where("teacher_id = ?", teacher.ID).First(&quiz).Error; err != nil {
		return err
	}

	qSubmission := &model.QuizSubmission{}
	if err := db.Where("quiz_id = ? AND id = ?", quiz.ID, id).Delete(&qSubmission).Error; err != nil {
		return err
	}

	return nil
}

func CreateQuizSubmission(db *gorm.DB, qSubmission *model.QuizSubmission) (*model.QuizSubmission, error) {
	if err := db.Create(qSubmission).Error; err != nil {
		return nil, err
	}

	return qSubmission, nil
}

func UpdateQuizSubmission(db *gorm.DB, qSubmission *model.QuizSubmission) error {
	if err := db.Model(&model.QuizSubmission{}).Where("quiz_id = ?", qSubmission.QuizID).Update(qSubmission).Error; err != nil {
		return err
	}

	return nil
}

func UpdateQuizSubmissionByOwner(db *gorm.DB, qSubmission *model.QuizSubmission, email string) error {
	owner := &model.User{}
	if err := db.Where("email = ?", email).First(&owner).Error; err != nil {
		return err
	}

	if err := db.Model(&model.QuizSubmission{}).Where("student_id = ? AND quiz_id = ?", owner.ID, qSubmission.QuizID).Update(qSubmission).Error; err != nil {
		return err
	}

	return nil
}

func UpdateQuizSubmissionByTeacher(db *gorm.DB, qSubmission *model.QuizSubmission, email string) error {
	teacher := &model.User{}
	if err := db.Where("email = ?", email).First(&teacher).Error; err != nil {
		return err
	}

	quiz := &model.Quiz{}
	if err := db.Where("teacher_id = ? AND id = ?", teacher.ID, qSubmission.QuizID).First(&quiz).Error; err != nil {
		return err
	}

	if err := db.Model(&model.QuizSubmission{}).Where("quiz_id = ?", quiz.ID).Update(qSubmission).Error; err != nil {
		return err
	}

	return nil
}
