package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"go.uber.org/zap"
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

	// 单机部署不需要处理跨域
	//handler := mws.CORS(mws.NewAuthBuilder().
	//	IgnorePath("/api/user/login").
	//	IgnorePath("/api/user/register").
	//	Auth(mux))
	handler := mws.NewAuthBuilder().
		IgnorePath("/api/user/login").
		IgnorePath("/api/user/register").
		Auth(mux)
	server := &http.Server{
		Addr:    "0.0.0.0:9091",
		Handler: handler,
	}

	log.Printf("HTTP server listening on %s", "localhost:9091")
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced shutting down err:%s\n", err)
	}

	log.Println("Server exited gracefully")
}

func initUserClient(client *clientv3.Client) user.UserServiceClient {
	cli := user.NewUserServiceClient(getSharedConn(client, userService))

	return cli
}

func initTaskClient(client *clientv3.Client) task.TaskServiceClient {
	cli := task.NewTaskServiceClient(getSharedConn(client, taskService))

	return cli
}

func getSharedConn(cli *clientv3.Client, serviceName string) *grpc.ClientConn {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	if conn, ok := connMap.Load(serviceName); ok {
		return conn.(*grpc.ClientConn)
	}

	svcCfg := `
	{
		"loadBalancingConfig": [
			{
				"round_robin": {}
			}
		]
	}`

	resolverBuilder, _ := resolver.NewBuilder(cli)
	conn, err := grpc.Dial(
		fmt.Sprintf("etcd:///%s", serviceName),
		grpc.WithResolvers(resolverBuilder),
		grpc.WithDefaultServiceConfig(svcCfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			logging.UnaryClientInterceptor(initInterceptor(logger))),
	)
	if err != nil {
		log.Fatal(err)
	}

	connMap.Store(serviceName, conn)
	return conn
}

func initInterceptor(l *zap.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, level logging.Level, msg string, fields ...any) {
		f := make([]zap.Field, 0, len(fields))

		for i := 0; i < len(fields); i += 2 {
			key := fields[i]
			value := fields[i+1]

			switch v := value.(type) {
			case string:
				f = append(f, zap.String(key.(string), v))
			case int:
				f = append(f, zap.Int(key.(string), v))
			case bool:
				f = append(f, zap.Bool(key.(string), v))
			default:
				f = append(f, zap.Any(key.(string), v))
			}
		}

		logger := l.WithOptions(zap.AddCallerSkip(1)).With(f...)

		switch level {
		case logging.LevelDebug:
			logger.Debug(msg)
		case logging.LevelInfo:
			logger.Info(msg)
		case logging.LevelWarn:
			logger.Warn(msg)
		case logging.LevelError:
			logger.Error(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", level))
		}
	})
}
