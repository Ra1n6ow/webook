package repository

import (
	"context"

	"github.com/ra1n6ow/webook/internal/domain"
	"github.com/ra1n6ow/webook/internal/repository/dao"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrRecordNotFound = dao.ErrRecordNotFound
)

type UserRepository struct {
	ud *dao.UserDAO
}

func NewUserRepository(ud *dao.UserDAO) *UserRepository {
	return &UserRepository{
		ud: ud,
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	ud := dao.User{
		Email:    u.Email,
		Password: u.Password,
	}
	return r.ud.Insert(ctx, ud)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.ud.QueryByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.toDomain(u), nil
}

func (r *UserRepository) Update(ctx context.Context, u domain.User) error {
	ud := dao.User{
		Id:       u.Id,
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		Intro:    u.Intro,
	}

	return r.ud.Update(ctx, ud)
}

func (r *UserRepository) FindById(ctx context.Context, userId int64) (domain.User, error) {
	u, err := r.ud.QueryById(ctx, userId)
	if err != nil {
		return domain.User{}, err
	}
	return r.toDomain(u), nil
}

func (r *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
		Birthday: u.Birthday,
		Nickname: u.Nickname,
		Intro:    u.Intro,
	}
}
