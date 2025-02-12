//go:build wireinject

package ioc

import (
	"fmt"
	"github.com/crazyfrankie/todolist/app/task/biz/repository"
	"github.com/crazyfrankie/todolist/app/task/biz/repository/dao"
	"github.com/crazyfrankie/todolist/app/task/biz/rpc"
	"github.com/crazyfrankie/todolist/app/task/biz/rpc/server"
	"github.com/crazyfrankie/todolist/app/task/biz/service"
	"github.com/crazyfrankie/todolist/app/task/config"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
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

	db.AutoMigrate(&dao.Task{})

	return db
}

func InitTask() *rpc.Server {
	wire.Build(
		InitDB,
		dao.NewTaskDao,
		repository.NewTaskRepo,
		service.NewTaskService,
		server.NewTaskServer,
		rpc.NewServer,
	)
	return new(rpc.Server)
}
