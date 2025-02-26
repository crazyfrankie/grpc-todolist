package service

import (
	"context"
	"errors"
	"strconv"

	"google.golang.org/grpc/metadata"

	"github.com/crazyfrankie/todolist/app/task/biz/repository"
	"github.com/crazyfrankie/todolist/app/task/biz/repository/dao"
)

type TaskService struct {
	repo *repository.TaskRepo
}

func NewTaskService(repo *repository.TaskRepo) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) AddTask(ctx context.Context, title, content string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("error param")
	}
	userId := md["user_id"][0]
	uId, _ := strconv.Atoi(userId)

	task := &dao.Task{
		UserId:  uId,
		Title:   title,
		Content: content,
	}

	return s.repo.CreateTask(ctx, task)
}

func (s *TaskService) List(ctx context.Context) ([]*dao.Task, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("error param")
	}
	userId := md["user_id"][0]
	uId, _ := strconv.Atoi(userId)

	return s.repo.FindByUid(ctx, uId, 0)
}

func (s *TaskService) UpdateTask(ctx context.Context, t *dao.Task) error {
	return s.repo.UpdateTask(ctx, t)
}

func (s *TaskService) DeleteTask(ctx context.Context, id int) error {
	return s.repo.DeleteTask(ctx, id)
}

func (s *TaskService) RecycleBin(ctx context.Context) ([]*dao.Task, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("error param")
	}
	userId := md["user_id"][0]
	uId, _ := strconv.Atoi(userId)

	return s.repo.FindByUid(ctx, uId, 1)
}

func (s *TaskService) RestoreTask(ctx context.Context, id int) error {
	return s.repo.RestoreTask(ctx, id)
}
