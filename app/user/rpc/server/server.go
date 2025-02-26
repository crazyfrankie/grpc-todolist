package server

import (
	"context"
	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"

	"github.com/crazyfrankie/todolist/app/user/biz/service"
	"github.com/crazyfrankie/todolist/app/user/rpc_gen/user"
)

type UserServer struct {
	svc *service.UserService
	user.UnimplementedUserServiceServer
}

func NewUserServer(svc *service.UserService) *UserServer {
	return &UserServer{svc: svc}
}

func (s *UserServer) RegisterServer(server *grpc.Server) {
	user.RegisterUserServiceServer(server, s)
}

func (s *UserServer) Register(ctx context.Context, request *user.RegisterRequest) (*user.RegisterResponse, error) {
	token, err := s.svc.Register(ctx, request.GetName(), request.GetPassword())
	if err != nil {
		return &user.RegisterResponse{
			Code: 200,
			Msg:  "ok",
		}, err
	}

	header := metadata.Pairs("Set-Auth-Token", token)
	grpc.SendHeader(ctx, header)

	return &user.RegisterResponse{
		Code: 200,
		Msg:  "ok",
	}, nil
}

func (s *UserServer) Login(ctx context.Context, request *user.LoginRequest) (*user.LoginResponse, error) {
	token, err := s.svc.Login(ctx, request.GetName(), request.GetPassword())
	if err != nil {
		return &user.LoginResponse{
			Code: 500,
			Msg:  "internal error",
		}, err
	}

	header := metadata.Pairs("Set-Auth-Token", token)
	grpc.SendHeader(ctx, header)

	return &user.LoginResponse{
		Code: 200,
		Msg:  "ok",
	}, nil
}

func (s *UserServer) GetUserInfo(ctx context.Context, request *user.GetUserInfoRequest) (*user.GetUserInfoResponse, error) {
	u, err := s.svc.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}

	newUser := &user.User{
		Id:   int32(u.Id),
		Name: u.Name,
	}

	return &user.GetUserInfoResponse{
		User: newUser,
	}, nil
}
