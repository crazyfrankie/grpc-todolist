package server

import (
	"context"
	"github.com/crazyfrankie/todolist/app/task/biz/repository/dao"
	"time"

	"google.golang.org/grpc"

	"github.com/crazyfrankie/todolist/app/task/biz/service"
	"github.com/crazyfrankie/todolist/app/task/rpc_gen/task"
)

type TaskServer struct {
	svc *service.TaskService
	task.UnimplementedTaskServiceServer
}

func NewTaskServer(svc *service.TaskService) *TaskServer {
	return &TaskServer{svc: svc}
}

func (t *TaskServer) RegisterServer(s *grpc.Server) {
	task.RegisterTaskServiceServer(s, t)
}

func (t *TaskServer) AddTask(ctx context.Context, request *task.AddTaskRequest) (*task.AddTaskResponse, error) {
	err := t.svc.AddTask(ctx, request.GetTitle(), request.GetContent())
	if err != nil {
		return nil, err
	}

	return &task.AddTaskResponse{}, nil
}

func (t *TaskServer) ListTasks(ctx context.Context, request *task.ListTasksRequest) (*task.ListTasksResponse, error) {
	tasks, err := t.svc.List(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]*task.Task, 0, len(tasks))
	for _, t := range tasks {
		utime := time.Unix(t.Utime, 0).Format(time.DateTime)
		results = append(results, &task.Task{
			Id:      int32(t.Id),
			Title:   t.Title,
			Content: t.Content,
			Utime:   utime,
			Status:  int32(t.Status),
		})
	}

	return &task.ListTasksResponse{
		Tasks: results,
	}, nil
}

func (t *TaskServer) UpdateTask(ctx context.Context, request *task.UpdateTaskRequest) (*task.UpdateTaskResponse, error) {
	tk := request.GetTask()
	err := t.svc.UpdateTask(ctx, &dao.Task{
		Id:      int(tk.GetId()),
		Title:   tk.GetTitle(),
		Content: tk.GetContent(),
	})
	if err != nil {
		return nil, err
	}

	return &task.UpdateTaskResponse{}, nil
}

func (t *TaskServer) DeleteTask(ctx context.Context, req *task.DeleteTaskRequest) (*task.DeleteTaskResponse, error) {
	err := t.svc.DeleteTask(ctx, int(req.GetId()))
	if err != nil {
		return nil, err
	}

	return &task.DeleteTaskResponse{}, nil
}

func (t *TaskServer) RecycleBin(ctx context.Context, req *task.RecycleBinRequest) (*task.RecycleBinResponse, error) {
	tasks, err := t.svc.RecycleBin(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]*task.Task, 0, len(tasks))
	for _, t := range tasks {
		utime := time.Unix(t.Utime, 0).Format(time.DateTime)
		results = append(results, &task.Task{
			Id:      int32(t.Id),
			Title:   t.Title,
			Content: t.Content,
			Utime:   utime,
			Status:  int32(t.Status),
		})
	}

	return &task.RecycleBinResponse{
		Tasks: results,
	}, nil
}
