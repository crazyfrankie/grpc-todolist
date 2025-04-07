//go:build wireinject

package ioc

import (
	"github.com/crazyfrankie/todolist/app/task/biz/repository"
	"github.com/crazyfrankie/todolist/app/task/biz/repository/dao"
	"github.com/crazyfrankie/todolist/app/task/biz/service"
	"github.com/crazyfrankie/todolist/app/task/rpc"
	"github.com/crazyfrankie/todolist/app/task/rpc/server"
	"github.com/crazyfrankie/todolist/app/task/rpc_gen/task"
	"github.com/google/wire"
	"google.golang.org/grpc"
)

func InitTask() *rpc.Server {
	wire.Build(
		InitDB,
		dao.NewTaskDao,
		repository.NewTaskRepo,
		service.NewTaskService,
		server.NewTaskServer,
		InitRegistry,
		registerService,
		rpc.NewServer,
	)
	return new(rpc.Server)
}

func registerService(t *server.TaskServer) func(grpc.ServiceRegistrar) {
	return func(s grpc.ServiceRegistrar) {
		task.RegisterTaskServiceServer(s, t)
	}
}
