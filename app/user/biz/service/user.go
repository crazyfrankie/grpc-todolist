package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"

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
	token, err = GenerateToken(u.Id)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) Login(ctx context.Context, name, password string) (string, error) {
	u, err := s.repo.FindByName(ctx, name)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return "", err
	}

	var token string
	token, err = GenerateToken(u.Id)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) GetUserInfo(ctx context.Context, id int) (dao.User, error) {
	return s.repo.FindById(ctx, id)
}
