package repository

import (
	"context"
	"github.com/crazyfrankie/todolist/app/user/biz/repository/dao"
)

type UserRepo struct {
	dao *dao.UserDao
}

func NewUserRepo(d *dao.UserDao) *UserRepo {
	return &UserRepo{dao: d}
}

func (r *UserRepo) CreateUser(ctx context.Context, u *dao.User) error {
	return r.dao.Create(ctx, u)
}

func (r *UserRepo) FindByName(ctx context.Context, name string) (dao.User, error) {
	return r.dao.FindByName(ctx, name)
}

func (r *UserRepo) FindById(ctx context.Context, id int) (dao.User, error) {
	return r.dao.FindById(ctx, id)
}
