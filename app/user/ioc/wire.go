//go:build wireinject

package ioc

import (
	"github.com/crazyfrankie/todolist/app/user/biz/repository"
	"github.com/crazyfrankie/todolist/app/user/biz/repository/dao"
	"github.com/crazyfrankie/todolist/app/user/biz/service"
	"github.com/crazyfrankie/todolist/app/user/rpc"
	"github.com/crazyfrankie/todolist/app/user/rpc/server"
	"github.com/crazyfrankie/todolist/app/user/rpc_gen/user"
	"github.com/google/wire"
	"google.golang.org/grpc"
)

func InitUser() *rpc.Server {
	wire.Build(
		InitDB,
		dao.NewUserDao,
		repository.NewUserRepo,
		service.NewUserService,
		server.NewUserServer,
		InitRegistry,
		registerService,
		rpc.NewServer,
	)
	return new(rpc.Server)
}

func registerService(u *server.UserServer) func(grpc.ServiceRegistrar) {
	return func(s grpc.ServiceRegistrar) {
		user.RegisterUserServiceServer(s, u)
	}
}
