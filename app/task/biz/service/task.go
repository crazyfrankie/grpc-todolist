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

func (s *TaskService) AddTask(ctx context.Context, title, content string, uid int) error {
	err := s.UserAuth(ctx, uid)
	if err != nil {
		return err
	}

	task := &dao.Task{
		UserId:  uid,
		Title:   title,
		Content: content,
	}

	return s.repo.CreateTask(ctx, task)
}

func (s *TaskService) List(ctx context.Context, uid int) ([]*dao.Task, error) {
	err := s.UserAuth(ctx, uid)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByUid(ctx, uid, 0)
}

func (s *TaskService) UpdateTask(ctx context.Context, t *dao.Task, uid int) error {
	err := s.UserAuth(ctx, uid)
	if err != nil {
		return err
	}
	return s.repo.UpdateTask(ctx, t)
}

func (s *TaskService) DeleteTask(ctx context.Context, id, uid int) error {
	err := s.UserAuth(ctx, uid)
	if err != nil {
		return err
	}
	return s.repo.DeleteTask(ctx, id)
}

func (s *TaskService) RecycleBin(ctx context.Context, uid int) ([]*dao.Task, error) {
	err := s.UserAuth(ctx, uid)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByUid(ctx, uid, 1)
}

func (s *TaskService) RestoreTask(ctx context.Context, id, uid int) error {
	err := s.UserAuth(ctx, uid)
	if err != nil {
		return err
	}
	return s.repo.RestoreTask(ctx, id)
}

func (s *TaskService) UserAuth(ctx context.Context, uid int) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("error param")
	}
	userId := md["user_id"][0]
	uId, _ := strconv.Atoi(userId)

	if uId != uid {
		return errors.New("UnAuthorized")
	}

	return nil
}
