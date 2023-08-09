package repository

import (
	"context"
	"time"

	"GeekTime/my-geektime/webook/internal/domain"
	"GeekTime/my-geektime/webook/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:        u.Id,
		Email:     u.Email,
		Password:  u.Password,
		Nickname:  u.Nickname,
		Biography: u.Biography,
		Birthday:  time.UnixMilli(u.Birthday).UTC(),
		Ctime:     time.UnixMilli(u.Ctime).UTC(),
	}, nil
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (r *UserRepository) Edit(ctx context.Context, id int, u domain.UserEdit) error {
	return r.dao.Update(ctx, id, dao.User{
		Email:     u.Email,
		Nickname:  u.Nickname,
		Biography: u.Biography,
		Birthday:  u.Birthday.UnixMilli(),
	})
}

func (r *UserRepository) FindById(ctx context.Context, id int) (domain.UserRead, error) {
	// 先从 cache 里面找
	// 再从 dao 里面找
	u, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.UserRead{}, err
	}
	// 找到了回写 cache
	return domain.UserRead{
		Email:     u.Email,
		Nickname:  u.Nickname,
		Biography: u.Biography,
		Birthday:  time.UnixMilli(u.Birthday).UTC(),
	}, nil
}
