package rpc

import (
	"context"
	"fmt"
	"github.com/crazyfrankie/todolist/app/user/rpc/server"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"log"
	"net"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"

	"github.com/crazyfrankie/todolist/app/user/config"
)

type Server struct {
	*grpc.Server
	Addr   string
	client *clientv3.Client
}

func NewServer(u *server.UserServer, client *clientv3.Client) *Server {
	s := grpc.NewServer()
	u.RegisterServer(s)

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
	em, err := endpoints.NewManager(cli, "service/user")
	if err != nil {
		return err
	}

	addr := "localhost" + port
	serviceKey := "service/user/" + addr

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

	return err
}
