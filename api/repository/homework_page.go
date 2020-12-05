package repository

import (
	"api/model"

	"github.com/jinzhu/gorm"
)

func ListHomeworkPagesByOwner(db *gorm.DB, email string) (model.HomeworkPages, error) {
	user := &model.User{}
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	homeworkPages := make([]*model.HomeworkPage, 0)
	if err := db.Where("teacher_id = ?", user.ID).Find(&homeworkPages).Error; err != nil {
		return nil, err
	}

	return homeworkPages, nil
}

func ReadHomeworkPage(db *gorm.DB, id uint) (*model.HomeworkPage, error) {
	homeworkPage := &model.HomeworkPage{}
	if err := db.Where("id = ?", id).First(&homeworkPage).Error; err != nil {
		return nil, err
	}

	return homeworkPage, nil
}

func ReadHomeworkPageWithNoOwner(db *gorm.DB, id uint) (*model.HomeworkPage, error) {
	homeworkPage := &model.HomeworkPage{}
	if err := db.Where("teacher_id is NULL AND id = ?", id).First(&homeworkPage).Error; err != nil {
		return nil, err
	}

	return homeworkPage, nil
}

func ReadHomeworkPageWithNoOwnerByStudentLink(db *gorm.DB, link string) (*model.HomeworkPage, error) {
	homeworkPage := &model.HomeworkPage{}
	if err := db.Where("teacher_id is NULL AND student_link = ?", link).First(&homeworkPage).Error; err != nil {
		return nil, err
	}

	return homeworkPage, nil
}

func ReadHomeworkPageWithNoOwnerByTeacherLink(db *gorm.DB, link string) (*model.HomeworkPage, error) {
	homeworkPage := &model.HomeworkPage{}
	if err := db.Where("teacher_id is NULL AND teacher_link = ?", link).First(&homeworkPage).Error; err != nil {
		return nil, err
	}

	return homeworkPage, nil
}

func ReadHomeworkPageByOwner(db *gorm.DB, id uint, email string) (*model.HomeworkPage, error) {
	user := &model.User{}
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	homeworkPage := &model.HomeworkPage{}
	if err := db.Where("teacher_id = ? AND id = ?", user.ID, id).First(&homeworkPage).Error; err != nil {
		return nil, err
	}

	return homeworkPage, nil
}

func ReadHomeworkPageWithOwnerByTeacherLink(db *gorm.DB, link, email string) (*model.HomeworkPage, error) {
	user := &model.User{}
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	homeworkPage := &model.HomeworkPage{}
	if err := db.Where("teacher_id = ? AND teacher_link = ?", user.ID, link).First(&homeworkPage).Error; err != nil {
		return nil, err
	}

	return homeworkPage, nil
}

func ReadHomeworkPageByOwnerByStudentLink(db *gorm.DB, link string) (*model.HomeworkPage, error) {
	homeworkPage := &model.HomeworkPage{}
	if err := db.Where("teacher_id IS NOT NULL AND student_link = ?", link).First(&homeworkPage).Error; err != nil {
		return nil, err
	}

	return homeworkPage, nil
}

func DeleteHomeworkPageByOwner(db *gorm.DB, id uint, email string) error {
	user := &model.User{}
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}

	homeworkPage := &model.HomeworkPage{}
	if err := db.Where("id = ? AND teacher_id = ?", id, user.ID).Delete(&homeworkPage).Error; err != nil {
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
	user := &model.User{}
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}

	if err := db.Where("teacher_id = ? AND id = ?", user.ID, homeworkPage.ID).First(&model.HomeworkPage{}).Update(homeworkPage).Error; err != nil {
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
