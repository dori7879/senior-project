package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Homeworks []*Homework

type Homework struct {
	gorm.Model
	Title     string
	Content   string
	LockingAt time.Time
}

type HomeworkDtos []*HomeworkDto

type HomeworkDto struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	LockingAt string `json:"locking_at"`
}

type HomeworkForm struct {
	Title     string `json:"title" form:"required, max=255"`
	Content   string `json:"content" form:"required, max=255"`
	LockingAt string `json:"locking_at" form:"required"`
}

func (hw Homework) ToDto() *HomeworkDto {
	return &HomeworkDto{
		ID:        hw.ID,
		Title:     hw.Title,
		Content:   hw.Content,
		LockingAt: hw.LockingAt.Format("2006-01-02"),
	}
}

func (hws Homeworks) ToDto() HomeworkDtos {
	dtos := make([]*HomeworkDto, len(hws))
	for i, hw := range hws {
		dtos[i] = hw.ToDto()
	}
	return dtos
}

func (f *HomeworkForm) ToModel() (*Homework, error) {
	lockingAt, err := time.Parse("2008-01-23", f.LockingAt)
	if err != nil {
		return nil, err
	}

	return &Homework{
		Title:     f.Title,
		Content:   f.Content,
		LockingAt: lockingAt,
	}, nil
}
