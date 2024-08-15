package task

import (
	"time"
)

type Task struct {
	ID       int       `gorm:"column:id"`
	Status   int       `gorm:"column:status"`
	Name     string    `gorm:"column:name"`
	CreateAt time.Time `gorm:"column:create_at"`
	FinishAt time.Time `gorm:"column:finish_at"`
}

func (Task) TableName() string {
	return "task"
}
