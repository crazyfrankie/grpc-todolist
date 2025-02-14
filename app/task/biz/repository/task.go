package repository

import (
	"context"

	"github.com/crazyfrankie/todolist/app/task/biz/repository/dao"
)

type TaskRepo struct {
	dao *dao.TaskDao
}

func NewTaskRepo(d *dao.TaskDao) *TaskRepo {
	return &TaskRepo{dao: d}
}

func (r *TaskRepo) CreateTask(ctx context.Context, t *dao.Task) error {
	return r.dao.Create(ctx, t)
}

func (r *TaskRepo) FindByUid(ctx context.Context, uid, status int) ([]*dao.Task, error) {
	return r.dao.FindByUid(ctx, uid, status)
}

func (r *TaskRepo) UpdateTask(ctx context.Context, t *dao.Task) error {
	return r.dao.UpdateTask(ctx, t)
}

func (r *TaskRepo) DeleteTask(ctx context.Context, id int) error {
	return r.dao.DeleteTask(ctx, id)
}

func (r *TaskRepo) RestoreTask(ctx context.Context, id int) error {
	return r.dao.RestoreTask(ctx, id)
}
