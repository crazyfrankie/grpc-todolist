package ioc

import (
	"time"

	"github.com/crazyfrankie/todolist/app/user/config"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func InitRegistry() *clientv3.Client {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:        []string{config.GetConf().ETCD.Addr},
		DialTimeout:      time.Second,
		AutoSyncInterval: 0,
	})
	if err != nil {
		panic(err)
	}

	return cli
}
