package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/carzyfrankie/app/gateway/mws"
	"github.com/crazyfrankie/todolist/app/task/rpc_gen/task"
	"github.com/crazyfrankie/todolist/app/user/rpc_gen/user"
)

var (
	userService = "service/user"
	taskService = "service/task"
	connMap     sync.Map
)

func main() {
	mux := runtime.NewServeMux(runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
		md := metadata.MD{}

		if userID, ok := request.Context().Value("user_id").(string); ok {
			md.Set("user_id", userID)
		}

		return md
	}))

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: time.Second * 5,
	})

	u := initUserClient(cli)
	t := initTaskClient(cli)

	err = user.RegisterUserServiceHandlerClient(context.Background(), mux, u)
	if err != nil {
		panic(err)
	}
	err = task.RegisterTaskServiceHandlerClient(context.Background(), mux, t)
	if err != nil {
		panic(err)
	}
	handler := mws.NewAuthBuilder().
		IgnorePath("/api/user/login").
		IgnorePath("/api/user/register").
		Auth(mux)

	log.Printf("HTTP server listening on %s", "localhost:9091")
	if err := http.ListenAndServe(fmt.Sprintf("%s", "localhost:9091"), handler); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func initUserClient(client *clientv3.Client) user.UserServiceClient {
	cli := user.NewUserServiceClient(getSharedConn(client, userService))

	go watchServices(client, userService)

	return cli
}

func initTaskClient(client *clientv3.Client) task.TaskServiceClient {
	cli := task.NewTaskServiceClient(getSharedConn(client, taskService))

	go watchServices(client, taskService)

	return cli
}

func getSharedConn(cli *clientv3.Client, serviceName string) *grpc.ClientConn {
	if conn, ok := connMap.Load(serviceName); ok {
		return conn.(*grpc.ClientConn)
	}

	resolverBuilder, _ := resolver.NewBuilder(cli)
	conn, err := grpc.Dial(
		fmt.Sprintf("etcd:///%s", serviceName),
		grpc.WithResolvers(resolverBuilder),
		grpc.WithDefaultServiceConfig(`{"LoadBalancingPolicy":"round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}

	connMap.Store(serviceName, conn)
	return conn
}

// 监听服务变更
func watchServices(cli *clientv3.Client, serviceName string) {
	watchChan := cli.Watch(context.Background(), serviceName, clientv3.WithPrefix())

	localServiceMap := &sync.Map{}

	for {
		select {
		case resp := <-watchChan:
			if resp.Err() != nil {
				log.Printf("Watch error: %v", resp.Err())
				continue
			}

			for _, ev := range resp.Events {
				key := string(ev.Kv.Key)
				addr := string(ev.Kv.Value)

				switch ev.Type {
				case clientv3.EventTypePut:
					if ev.IsCreate() {
						log.Printf("New service registered: %s", addr)
						localServiceMap.Store(key, addr)
					} else if ev.IsModify() {
						log.Printf("Service updated: %s", addr)
						localServiceMap.Store(key, addr)
					}
				case clientv3.EventTypeDelete:
					log.Printf("Service unregistered: %s", addr)
					localServiceMap.Delete(key)
				}

				var services []string
				localServiceMap.Range(func(k, v interface{}) bool {
					services = append(services, v.(string))
					return true
				})
				log.Printf("Current services: %v", services)
			}
		}
	}
}
