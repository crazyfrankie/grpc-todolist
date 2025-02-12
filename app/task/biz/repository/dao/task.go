package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Task struct {
	Id      int    `gorm:"primaryKey,autoIncrement"`
	UserId  int    `gorm:"index:user_utime"`
	Title   string `gorm:"type:varchar(128)"`
	Content string
	Status  int
	Ctime   int64
	Utime   int64 `gorm:"index:user_utime"`
}

type TaskDao struct {
	db *gorm.DB
}

func NewTaskDao(db *gorm.DB) *TaskDao {
	return &TaskDao{db: db}
}

func (d *TaskDao) Create(ctx context.Context, t *Task) error {
	now := time.Now().Unix()
	t.Ctime = now
	t.Utime = now

	return d.db.WithContext(ctx).Create(t).Error
}

func (d *TaskDao) FindById(ctx context.Context, id int) (Task, error) {
	var t Task
	err := d.db.WithContext(ctx).Model(&Task{}).Where("id = ?", id).Find(&t).Error
	if err != nil {
		return Task{}, err
	}

	return t, nil
}

func (d *TaskDao) FindByUid(ctx context.Context, uid, status int) ([]*Task, error) {
	var tasks []*Task
	err := d.db.WithContext(ctx).Model(&Task{}).Where("user_id = ? AND status = ?", uid, status).Order("utime DESC").Find(&tasks).Error
	if err != nil {
		return []*Task{}, err
	}

	return tasks, nil
}

func (d *TaskDao) UpdateTask(ctx context.Context, t *Task) error {
	return d.db.WithContext(ctx).Model(&Task{}).Where("id = ?", t.Id).Updates(t).Error
}

func (d *TaskDao) DeleteTask(ctx context.Context, id int) error {
	return d.db.WithContext(ctx).Model(&Task{}).Where("id = ?", id).UpdateColumn("status", 1).Error
}
