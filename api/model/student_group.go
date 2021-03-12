package model

type StudentGroups []*StudentGroup

type StudentGroup struct {
	ID        uint `gorm:"primaryKey"`
	Title     string
	ShareLink string
	OwnerID   uint      `gorm:"default:null"`
	Teachers  []Teacher `gorm:"many2many:teachers_student_groups;"`
	Students  []Student `gorm:"many2many:students_student_groups;"`
}

type StudentGroupDtos []*StudentGroupDto

type StudentGroupDto struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	ShareLink string `json:"share_link"`
	OwnerID   uint   `json:"owner_id"`
}

type StudentGroupForm struct {
	Title     string `json:"title" form:"max=255"`
	ShareLink string `json:"share_link" form:""`
	OwnerID   uint   `json:"owner_id" form:""`
}

func (ac StudentGroup) ToDto() *StudentGroupDto {
	return &StudentGroupDto{
		ID:        ac.ID,
		Title:     ac.Title,
		ShareLink: ac.ShareLink,
		OwnerID:   ac.OwnerID,
	}
}

func (acs StudentGroups) ToDto() StudentGroupDtos {
	dtos := make([]*StudentGroupDto, len(acs))
	for i, ac := range acs {
		dtos[i] = ac.ToDto()
	}
	return dtos
}

func (f *StudentGroupForm) ToModel() (*StudentGroup, error) {
	return &StudentGroup{
		Title:     f.Title,
		ShareLink: f.ShareLink,
		OwnerID:   f.OwnerID,
	}, nil
}
