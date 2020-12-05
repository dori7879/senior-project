package repository

import (
	"api/model"

	"github.com/jinzhu/gorm"
)

func ListHomeworksByOwnerEmail(db *gorm.DB, email string) (model.Homeworks, error) {
	owner := &model.Student{}
	if err := db.Where("email = ?", email).First(&owner).Error; err != nil {
		return nil, err
	}

	homeworks := make([]*model.Homework, 0)
	if err := db.Where("student_id = ?", owner.ID).Find(&homeworks).Error; err != nil {
		return nil, err
	}

	return homeworks, nil
}

func ListRelatedHomeworks(db *gorm.DB, hwpID uint) (model.Homeworks, error) {
	homeworks := make([]*model.Homework, 0)
	if err := db.Where("homework_page_id = ?", hwpID).Find(&homeworks).Error; err != nil {
		return nil, err
	}

	return homeworks, nil
}

func ReadHomeworkByIDandOwner(db *gorm.DB, id uint, email string) (*model.Homework, error) {
	owner := &model.User{}
	if err := db.Where("email = ?", email).First(&owner).Error; err != nil {
		return nil, err
	}

	homework := &model.Homework{}
	if err := db.Where("student_id = ? AND id = ?", owner.ID, id).First(&homework).Error; err != nil {
		return nil, err
	}

	return homework, nil
}

func ReadHomeworkByIDandTeacher(db *gorm.DB, id uint, email string) (*model.Homework, error) {
	teacher := &model.User{}
	if err := db.Where("email = ?", email).First(&teacher).Error; err != nil {
		return nil, err
	}

	hwp := &model.HomeworkPage{}
	if err := db.Where("teacher_id = ?", teacher.ID).First(&hwp).Error; err != nil {
		return nil, err
	}

	homework := &model.Homework{}
	if err := db.Where("id = ? AND homework_page_id = ?", id, hwp.ID).First(&homework).Error; err != nil {
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

func DeleteHomeworkByTeacher(db *gorm.DB, id uint, email string) error {
	teacher := &model.User{}
	if err := db.Where("email = ?", email).First(&teacher).Error; err != nil {
		return err
	}

	hwp := &model.HomeworkPage{}
	if err := db.Where("teacher_id = ?", teacher.ID).First(&hwp).Error; err != nil {
		return err
	}

	homework := &model.Homework{}
	if err := db.Where("homework_page_id = ? AND id = ?", hwp.ID, id).Delete(&homework).Error; err != nil {
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
	if err := db.Model(&model.Homework{}).Where("homework_page_id = ?", homework.HomeworkPageID).Update(homework).Error; err != nil {
		return err
	}

	return nil
}

func UpdateHomeworkByOwner(db *gorm.DB, homework *model.Homework, email string) error {
	owner := &model.User{}
	if err := db.Where("email = ?", email).First(&owner).Error; err != nil {
		return err
	}

	if err := db.Model(&model.Homework{}).Where("student_id = ? AND homework_page_id = ?", owner.ID, homework.HomeworkPageID).Update(homework).Error; err != nil {
		return err
	}

	return nil
}

func UpdateHomeworkByTeacher(db *gorm.DB, homework *model.Homework, email string) error {
	teacher := &model.User{}
	if err := db.Where("email = ?", email).First(&teacher).Error; err != nil {
		return err
	}

	hwp := &model.HomeworkPage{}
	if err := db.Where("teacher_id = ? AND id = ?", teacher.ID, homework.HomeworkPageID).First(&hwp).Error; err != nil {
		return err
	}

	if err := db.Model(&model.Homework{}).Where("homework_page_id = ?", hwp.ID).Update(homework).Error; err != nil {
		return err
	}

	return nil
}
