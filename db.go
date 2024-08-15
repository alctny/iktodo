package main

import (
	"errors"

	"github.com/urfave/cli/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbFile = "/home/ant/.config/iktodo.db"

var db *gorm.DB

func GetDBFile() string {
	return dbFile
}

func initDB(*cli.Context) error {
	var err error
	dbFile := GetDBFile()
	db, err = gorm.Open(sqlite.Open(dbFile))
	return err
}

// SaveTask 新建任务
func SaveTask(t *Task) error {
	return db.Create(t).Error
}

// DeleteTask 删除任务
func DeleteTask(id int) error {
	return db.Delete(Task{ID: id}).Error
}

// ListTask 查询任务列表
func ListTask(all bool) ([]Task, error) {
	var lis []Task
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
	var t Task
	tx := db.Where(Task{ID: id}).Take(&t)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return errors.New("task id not exist")
		}
		return tx.Error
	}

	status := ^t.Status
	return db.Model(Task{}).Where(Task{ID: id}).Update("status", status).Error
}
