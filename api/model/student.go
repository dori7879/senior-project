package model

type Students []*Student

type Student struct {
	ID        uint `gorm:"column:student_id"`
	Homeworks []Homework
}

type StudentDtos []*StudentDto

type StudentDto struct {
	ID        uint       `json:"id"`
	Homeworks []Homework `json:"homeworks"`
}

func (s Student) ToDto() *StudentDto {
	return &StudentDto{
		ID:        s.ID,
		Homeworks: s.Homeworks,
	}
}

func (ss Students) ToDto() StudentDtos {
	dtos := make([]*StudentDto, len(ss))
	for i, s := range ss {
		dtos[i] = s.ToDto()
	}
	return dtos
}
