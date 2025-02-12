package rpc

import (
	"net"

	"google.golang.org/grpc"

	"github.com/crazyfrankie/todolist/app/user/biz/rpc/server"
	"github.com/crazyfrankie/todolist/app/user/config"
)

type Server struct {
	*grpc.Server
	Addr string
}

func NewServer(u *server.UserServer) *Server {
	s := grpc.NewServer()
	u.RegisterServer(s)

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
