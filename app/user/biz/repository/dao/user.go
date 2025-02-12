package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id       int    `gorm:"primaryKey,autoIncrement"`
	Name     string `gorm:"unique;type:varchar(255)"`
	Password string `gorm:";type:varchar(255)"`
	Ctime    int64
	Utime    int64
}
type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

func (d *UserDao) Create(ctx context.Context, u *User) error {
	now := time.Now().Unix()
	u.Ctime = now
	u.Utime = now

	return d.db.WithContext(ctx).Create(u).Error
}

func (d *UserDao) FindByName(ctx context.Context, name string) (User, error) {
	var u User
	err := d.db.WithContext(ctx).Model(&User{}).Where("name = ?", name).Find(&u).Error
	if err != nil {
		return User{}, err
	}

	return u, nil
}

func (d *UserDao) FindById(ctx context.Context, id int) (User, error) {
	var u User
	err := d.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Find(&u).Error
	if err != nil {
		return User{}, err
	}

	return u, nil
}
