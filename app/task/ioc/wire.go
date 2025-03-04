//go:build wireinject

package ioc

import (
	"github.com/crazyfrankie/todolist/app/task/biz/repository"
	"github.com/crazyfrankie/todolist/app/task/biz/repository/dao"
	"github.com/crazyfrankie/todolist/app/task/biz/service"
	"github.com/crazyfrankie/todolist/app/task/rpc"
	"github.com/crazyfrankie/todolist/app/task/rpc/server"
	"github.com/google/wire"
)

func InitTask() *rpc.Server {
	wire.Build(
		InitDB,
		dao.NewTaskDao,
		repository.NewTaskRepo,
		service.NewTaskService,
		server.NewTaskServer,
		InitRegistry,
		rpc.NewServer,
	)
	return new(rpc.Server)
}
