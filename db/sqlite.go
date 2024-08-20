package db

import (
	"fmt"
	"os"
	"path"

	"github.com/alctny/iktodo/common"
	"github.com/alctny/iktodo/task"
	"github.com/urfave/cli/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbFile = ".iktodo.db"

var db *gorm.DB

func GetDBFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "panic: get home path error: %v", err)
		os.Exit(1)
	}
	sqlPath := path.Join(home, dbFile)
	return sqlPath
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
func DeleteTask(id []int) error {
	return db.Model(task.Task{}).Where("id IN (?)", id).Delete(nil).Error
}

// ListTask 查询任务列表
func ListTask(query map[string]any, limit uint, offset uint) ([]task.Task, error) {
	var lis []task.Task
	tx := db.Where(query).Order("status DESC, create_at")
	if limit != 0 {
		tx = tx.Offset(int(offset)).Limit(int(limit))
	}
	err := tx.Find(&lis).Error
	return lis, err
}

// DoneTask 将任务标识为完成/未完成
func DoneTask(id []int) error {
	var ts []task.Task
	err := db.Where("id IN (?)", id).Find(&ts).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return common.ErrTaskId
		}
		return err
	}

	// TODO 使用事务
	for _, t := range ts {
		status := ^t.Status
		err = db.Model(task.Task{}).Where(task.Task{ID: t.ID}).Update("status", status).Error
		if err != nil {
			return err
		}
	}
	return nil
}

// RemoveFinished 删除所有已完成任务
func RemoveFinished() error {
	return db.Model(task.Task{}).Where(task.Task{Status: -1}).Delete(nil).Error
}

func Aggregate() (*task.AggregateResult, error) {
	var total int64
	err := db.Model(task.Task{}).Count(&total).Error
	if err != nil {
		return nil, err
	}

	var finished int64
	err = db.Model(task.Task{}).Where(map[string]any{"status": task.StatusFinished}).Count(&finished).Error
	if err != nil {
		return nil, err
	}

	var unfinish int64
	err = db.Model(task.Task{}).Where(map[string]any{"status": task.StatusUnfinish}).Count(&unfinish).Error
	if err != nil {
		return nil, err
	}

	return &task.AggregateResult{Total: uint(total), Finished: uint(finished), Unfinish: uint(unfinish)}, nil
}

// Search 搜索任务
func Search(kw string) ([]task.Task, error) {
	var ts []task.Task
	err := db.Model(task.Task{}).Where("name LIKE ?", "%"+kw+"%").Find(&ts).Error
	return ts, err
}

// UpdateTask 更新任务
func Update(id int, m map[string]any) error {
	return db.Model(task.Task{}).Where("id = ?", id).Updates(m).Error
}
