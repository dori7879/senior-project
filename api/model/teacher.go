package model

type Teachers []*Teacher

type Teacher struct {
	ID            uint `gorm:"column:teacher_id"`
	HomeworkPages []HomeworkPage
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
