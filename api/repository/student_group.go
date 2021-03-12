package repository

import (
	"api/model"

	"github.com/jinzhu/gorm"
)

func ListStudentGroupsByOwner(db *gorm.DB, email string) (model.StudentGroups, error) {
	owner := &model.User{}
	if err := db.Where("email = ?", email).First(&owner).Error; err != nil {
		return nil, err
	}

	sgs := make([]*model.StudentGroup, 0)
	if err := db.Where("owner_id = ?", owner.ID).Find(&sgs).Error; err != nil {
		return nil, err
	}

	return sgs, nil
}

func ListStudentGroupsByStudent(db *gorm.DB, email string) (model.StudentGroups, error) {
	user := &model.User{}
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	student := &model.Student{}
	sgs := make([]*model.StudentGroup, 0)
	if err := db.Model(&student).Where("student_id = ?", user.ID).Association("StudentGroups").Find(&sgs).Error; err != nil {
		return nil, err
	}

	return sgs, nil
}

func ListStudentGroupsByTeacher(db *gorm.DB, email string) (model.StudentGroups, error) {
	user := &model.User{}
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	teacher := &model.Teacher{}
	sgs := make([]*model.StudentGroup, 0)
	if err := db.Model(&teacher).Where("teacher_id = ?", user.ID).Association("StudentGroups").Find(&sgs).Error; err != nil {
		return nil, err
	}

	return sgs, nil
}

func ReadStudentGroup(db *gorm.DB, id uint) (*model.StudentGroup, error) {
	sg := &model.StudentGroup{}
	if err := db.Where("id = ?", id).First(&sg).Error; err != nil {
		return nil, err
	}

	return sg, nil
}

func DeleteStudentGroupByOwner(db *gorm.DB, id uint, email string) error {
	user := &model.User{}
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}

	sg := &model.StudentGroup{}
	if err := db.Where("id = ? AND owner_id = ?", id, user.ID).Delete(&sg).Error; err != nil {
		return err
	}

	return nil
}

func CreateStudentGroup(db *gorm.DB, sg *model.StudentGroup) (*model.StudentGroup, error) {
	if err := db.Create(sg).Error; err != nil {
		return nil, err
	}

	return sg, nil
}

func UpdateStudentGroupByOwner(db *gorm.DB, sg *model.StudentGroup, email string) error {
	user := &model.User{}
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}

	if err := db.Where("owner_id = ? AND id = ?", user.ID, sg.ID).First(&model.StudentGroup{}).Update(sg).Error; err != nil {
		return err
	}

	return nil
}

// func ListStudentsOfStudentGroup(db *gorm.DB, id uint) (model.Users, error) {
// 	sg := &model.StudentGroup{}
// 	students := make([]*model.User, 0)
// 	if err := db.Model(&sg).Where("id = ?", id).Association("Students").Find(&students).Error; err != nil {
// 		return nil, err
// 	}

// 	return sgs, nil
// }

// func AddStudentToStudentGroup(db *gorm.DB, id uint, email string) error {
// 	return nil
// }

// func AddTeacherToStudentGroup(db *gorm.DB, id uint, email string) error {
// 	return nil
// }
