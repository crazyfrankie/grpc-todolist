package rpc

import (
	"net"

	"google.golang.org/grpc"

	"github.com/crazyfrankie/todolist/app/task/biz/rpc/server"
	"github.com/crazyfrankie/todolist/app/task/config"
)

type Server struct {
	*grpc.Server
	Addr string
}

func NewServer(t *server.TaskServer) *Server {
	s := grpc.NewServer()
	t.RegisterServer(s)

	return &Server{
		Server: s,
		Addr:   config.GetConf().Server.Addr,
	}
}

func (s *Server) Serve() error {
	conn, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	return s.Server.Serve(conn)
}
