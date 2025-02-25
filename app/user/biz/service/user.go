package service

import (
	"context"
	"errors"
	"strconv"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/metadata"

	"github.com/crazyfrankie/todolist/app/user/biz/repository"
	"github.com/crazyfrankie/todolist/app/user/biz/repository/dao"
)

type UserService struct {
	repo *repository.UserRepo
}

func NewUserService(r *repository.UserRepo) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) Register(ctx context.Context, name, password string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("error param")
	}
	userAgent := md["user_agent"][0]

	ps, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	u := &dao.User{
		Name:     name,
		Password: string(ps),
	}

	err = s.repo.CreateUser(ctx, u)
	if err != nil {
		return "", err
	}

	var token string
	token, err = GenerateToken(u.Id, userAgent)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) Login(ctx context.Context, name, password string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("error param")
	}
	userAgent := md["user_agent"][0]

	u, err := s.repo.FindByName(ctx, name)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return "", err
	}

	var token string
	token, err = GenerateToken(u.Id, userAgent)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) GetUserInfo(ctx context.Context) (dao.User, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return dao.User{}, errors.New("error param")
	}
	userId := md["user_id"][0]
	uId, _ := strconv.Atoi(userId)

	return s.repo.FindById(ctx, uId)
}
