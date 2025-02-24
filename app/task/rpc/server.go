package rpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/crazyfrankie/todolist/app/task/rpc/server"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/crazyfrankie/todolist/app/task/config"
)

type Server struct {
	*grpc.Server
	Addr   string
	client *clientv3.Client
}

func NewServer(t *server.TaskServer, client *clientv3.Client) *Server {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer(grpc.ChainUnaryInterceptor(
		logging.UnaryServerInterceptor(initInterceptor(logger))),
	)
	t.RegisterServer(s)

	return &Server{
		Server: s,
		Addr:   config.GetConf().Server.Addr,
		client: client,
	}
}

func (s *Server) Serve() error {
	conn, err := net.Listen("tcp", s.Addr)
	if err != nil {
		panic(err)
	}

	err = registerServer(s.client, s.Addr)
	if err != nil {
		return err
	}

	return s.Server.Serve(conn)
}

func (s *Server) ShutDown() error {
	err := s.client.Close()
	if err != nil {
		return err
	}

	s.Server.GracefulStop()

	return nil
}

func registerServer(cli *clientv3.Client, port string) error {
	em, err := endpoints.NewManager(cli, "service/task")
	if err != nil {
		return err
	}

	addr := "localhost" + port
	serviceKey := "service/task/" + addr

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	leaseResp, err := cli.Grant(ctx, 60)
	if err != nil {
		log.Fatalf("failed to grant lease: %v", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = em.AddEndpoint(ctx, serviceKey, endpoints.Endpoint{Addr: addr}, clientv3.WithLease(leaseResp.ID))
	//_, err = cli.Put(ctx, serviceKey, addr, clientv3.WithLease(leaseResp.ID))

	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		ch, err := cli.KeepAlive(ctx, leaseResp.ID)
		if err != nil {
			log.Fatalf("KeepAlive failed: %v", err)
		}

		for {
			select {
			case _, ok := <-ch:
				if !ok { // 通道关闭，租约停止
					log.Println("KeepAlive channel closed")
					return
				}
				fmt.Println("Lease renewed")
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
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
			case bool:
				f = append(f, zap.Bool(key.(string), v))
			case int:
				f = append(f, zap.Int(key.(string), v))
			default:
				f = append(f, zap.Any(key.(string), v))
			}
		}

		logger := l.WithOptions(zap.AddCallerSkip(1)).With(f...)

		switch level {
		case logging.LevelInfo:
			logger.Info(msg)
		case logging.LevelDebug:
			logger.Debug(msg)
		case logging.LevelError:
			logger.Error(msg)
		case logging.LevelWarn:
			logger.Warn(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", level))
		}
	})
}
