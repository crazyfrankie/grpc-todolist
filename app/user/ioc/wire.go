//go:build wireinject

package ioc

import (
	"github.com/crazyfrankie/todolist/app/user/biz/repository"
	"github.com/crazyfrankie/todolist/app/user/biz/repository/dao"
	"github.com/crazyfrankie/todolist/app/user/biz/service"
	"github.com/crazyfrankie/todolist/app/user/rpc"
	"github.com/crazyfrankie/todolist/app/user/rpc/server"
	"github.com/google/wire"
)

func InitUser() *rpc.Server {
	wire.Build(
		InitDB,
		dao.NewUserDao,
		repository.NewUserRepo,
		service.NewUserService,
		server.NewUserServer,
		InitRegistry,
		rpc.NewServer,
	)
	return new(rpc.Server)
}
