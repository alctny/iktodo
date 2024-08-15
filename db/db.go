package db

import (
	"errors"

	"github.com/alctny/iktodo/task"
	"github.com/urfave/cli/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbFile = "/home/ant/.config/iktodo.db"

var db *gorm.DB

func GetDBFile() string {
	return dbFile
}

func InitDB(*cli.Context) error {
	var err error
	dbFile := GetDBFile()
	db, err = gorm.Open(sqlite.Open(dbFile))
	return err
}

// SaveTask 新建任务
func SaveTask(t *task.Task) error {
	return db.Create(t).Error
}

// DeleteTask 删除任务
func DeleteTask(id int) error {
	return db.Delete(task.Task{ID: id}).Error
}

// ListTask 查询任务列表
func ListTask(all bool) ([]task.Task, error) {
	var lis []task.Task
	tx := db.Session(&gorm.Session{})
	if !all {
		tx = db.Where(map[string]any{"status": 0})
	}
	tx = tx.Order("create_at").
		Find(&lis)
	return lis, tx.Error
}

// DoneTask 将任务标识为完成/未完成
func DoneTask(id int) error {
	var t task.Task
	tx := db.Where(task.Task{ID: id}).Take(&t)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return errors.New("task id not exist")
		}
		return tx.Error
	}

	status := ^t.Status
	return db.Model(task.Task{}).Where(task.Task{ID: id}).Update("status", status).Error
}
