package model

type Teachers []*Teacher

type Teacher struct {
	ID                 uint           `gorm:"primaryKey;column:teacher_id"`
	User               User           `gorm:"foreignKey:ID;references:ID"`
	HomeworkPages      []HomeworkPage `gorm:"foreignKey:TeacherID;references:ID"`
	OwnedStudentGroups []StudentGroup `gorm:"foreignKey:OwnerID;references:ID"`
	StudentGroups      []StudentGroup `gorm:"many2many:teachers_student_groups;"`
}

type TeacherDtos []*TeacherDto

type TeacherDto struct {
	ID            uint           `json:"id"`
	HomeworkPages []HomeworkPage `json:"homework_pages"`
}

func (t Teacher) ToDto() *TeacherDto {
	return &TeacherDto{
		ID:            t.ID,
		HomeworkPages: t.HomeworkPages,
	}
}

func (ts Teachers) ToDto() TeacherDtos {
	dtos := make([]*TeacherDto, len(ts))
	for i, t := range ts {
		dtos[i] = t.ToDto()
	}
	return dtos
}
