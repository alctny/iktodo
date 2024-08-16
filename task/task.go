package task

import (
	"fmt"
	"time"
)

type ColorWhen string

const (
	WhenNone ColorWhen = "none"
	WhenDone ColorWhen = "done"
	WhenAny  ColorWhen = "all"
)

const (
	StatusFinished = -1
	StatusUnfinish = 0
)

const (
	fmtNocolor  = "%02d  %s"
	fmtFinished = "%02d  \033[0;93;32m%s\033[0m"
	fmtUnfinish = "%02d  \033[1;91;35m%s\033[0m"
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

// ColorString 获取彩色输出字符串（不含换行符）
// TODO: 这个输出输出层，展示层，应该移动到其他磨快当中
func (t Task) ColorString(when ColorWhen) string {
	switch when {
	case WhenNone:
		return fmt.Sprintf(fmtNocolor, t.ID, t.Name)
	case WhenDone:
		if t.Status == StatusFinished {
			return fmt.Sprintf(fmtFinished, t.ID, t.Name)
		}
		return fmt.Sprintf(fmtUnfinish, t.ID, t.Name)
	case WhenAny:
		if t.Status == StatusFinished {
			return fmt.Sprintf(fmtFinished, t.ID, t.Name)
		}
		return fmt.Sprintf(fmtUnfinish, t.ID, t.Name)
	default:
		return t.Name
	}
}
