//go:build wireinject

package ioc

import (
	"fmt"
	"github.com/crazyfrankie/todolist/app/user/biz/repository"
	"github.com/crazyfrankie/todolist/app/user/biz/repository/dao"
	"github.com/crazyfrankie/todolist/app/user/biz/rpc/server"
	"github.com/crazyfrankie/todolist/app/user/biz/service"
	"github.com/crazyfrankie/todolist/app/user/config"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"

	"github.com/crazyfrankie/todolist/app/user/biz/rpc"
)

func InitDB() *gorm.DB {
	dsn := fmt.Sprintf(config.GetConf().MySQL.DSN,
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DB"))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: &schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&dao.User{})

	return db
}

func InitUser() *rpc.Server {
	wire.Build(
		InitDB,
		dao.NewUserDao,
		repository.NewUserRepo,
		service.NewUserService,
		server.NewUserServer,
		rpc.NewServer,
	)
	return new(rpc.Server)
}
