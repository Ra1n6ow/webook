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

func (r *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Email:    u.Email,
		Password: u.Password,
	}
}
