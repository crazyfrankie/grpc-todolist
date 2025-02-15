//go:build wireinject

package ioc

import (
	"fmt"
	"github.com/crazyfrankie/todolist/app/task/biz/repository"
	"github.com/crazyfrankie/todolist/app/task/biz/repository/dao"
	"github.com/crazyfrankie/todolist/app/task/biz/service"
	"github.com/crazyfrankie/todolist/app/task/config"
	"github.com/crazyfrankie/todolist/app/task/rpc"
	"github.com/crazyfrankie/todolist/app/task/rpc/server"
	"github.com/google/wire"
	clientv3 "go.etcd.io/etcd/client/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
	"time"
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

func InitRegistry() *clientv3.Client {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{config.GetConf().ETCD.Addr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	return cli
}

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
