package repository

import (
	"api/model"

	"github.com/jinzhu/gorm"
)

func ListHomeworkPagesByOwner(db *gorm.DB, email string) (model.HomeworkPages, error) {
	teacher := &model.Teacher{}
	if err := db.Where("email = ?", email).First(&teacher).Error; err != nil {
		return nil, err
	}

	homeworkPages := make([]*model.HomeworkPage, 0)
	if err := db.Where("teacher_id = ?", teacher.ID).Find(&homeworkPages).Error; err != nil {
		return nil, err
	}

	return homeworkPages, nil
}

func ReadHomeworkPageWithNoOwner(db *gorm.DB, id uint) (*model.HomeworkPage, error) {
	homeworkPage := &model.HomeworkPage{}
	if err := db.Where("teacher_id is NULL AND id = ?", id).First(&homeworkPage).Error; err != nil {
		return nil, err
	}

	return homeworkPage, nil
}

func ReadHomeworkPageByOwner(db *gorm.DB, id uint, email string) (*model.HomeworkPage, error) {
	teacher := &model.Teacher{}
	if err := db.Where("email = ?", email).First(&teacher).Error; err != nil {
		return nil, err
	}

	homeworkPage := &model.HomeworkPage{}
	if err := db.Where("teacher_id = ? AND id = ?", teacher.ID, id).First(&homeworkPage).Error; err != nil {
		return nil, err
	}

	return homeworkPage, nil
}

func DeleteHomeworkPageByOwner(db *gorm.DB, id uint, email string) error {
	teacher := &model.Teacher{}
	if err := db.Where("email = ?", email).First(&teacher).Error; err != nil {
		return err
	}

	homeworkPage := &model.HomeworkPage{}
	if err := db.Where("id = ? AND teacher_id = ?", id, teacher.ID).Delete(&homeworkPage).Error; err != nil {
		return err
	}

	return nil
}

func CreateHomeworkPage(db *gorm.DB, homeworkPage *model.HomeworkPage) (*model.HomeworkPage, error) {
	if err := db.Create(homeworkPage).Error; err != nil {
		return nil, err
	}

	return homeworkPage, nil
}

func UpdateHomeworkPage(db *gorm.DB, homeworkPage *model.HomeworkPage) error {
	if err := db.First(&model.HomeworkPage{}, homeworkPage.ID).Update(homeworkPage).Error; err != nil {
		return err
	}

	return nil
}

func UpdateHomeworkPageByOwner(db *gorm.DB, homeworkPage *model.HomeworkPage, email string) error {
	teacher := &model.Teacher{}
	if err := db.Where("email = ?", email).First(&teacher).Error; err != nil {
		return err
	}

	if err := db.Where("teacher_id = ? AND id = ?", teacher.ID, homeworkPage.ID).First(&model.HomeworkPage{}).Update(homeworkPage).Error; err != nil {
		return err
	}

	return nil
}

func UpdateHomeworkPageWithNoOwner(db *gorm.DB, homeworkPage *model.HomeworkPage) error {
	if err := db.Where("teacher_id is NULL AND id = ?", homeworkPage.ID).First(&model.HomeworkPage{}).Update(homeworkPage).Error; err != nil {
		return err
	}

	return nil
}
