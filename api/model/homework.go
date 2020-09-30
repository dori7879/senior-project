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
