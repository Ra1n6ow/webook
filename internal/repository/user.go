package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/ra1n6ow/webook/internal/domain"
	"github.com/ra1n6ow/webook/internal/repository/cache"
	"github.com/ra1n6ow/webook/internal/repository/dao"
)

var (
	ErrDuplicateUser = dao.ErrDuplicateEmail
	ErrUserNotFound  = dao.ErrRecordNotFound
)

type UserRepositoryer interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	Update(ctx context.Context, u domain.User) error
	FindById(ctx context.Context, userId int64) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
}

type UserRepository struct {
	dao   dao.UserDAOer
	cache cache.UserCacher
}

func NewUserRepository(dao dao.UserDAOer, cache cache.UserCacher) UserRepositoryer {
	return &UserRepository{
		dao:   dao,
		cache: cache,
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.toEntity(u))
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.QueryByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.toDomain(u), nil
}

func (r *UserRepository) Update(ctx context.Context, u domain.User) error {
	return r.dao.Update(ctx, r.toEntity(u))
}

func (r *UserRepository) FindById(ctx context.Context, userId int64) (domain.User, error) {
	uc, err := r.cache.Get(ctx, userId)
	if err == nil {
		return uc, nil
	}
	// 缓存未命中，从数据库中查询
	ud, err := r.dao.QueryById(ctx, userId)
	if err != nil {
		return domain.User{}, err
	}
	// 将数据库中的数据写入缓存
	u := r.toDomain(ud)
	err = r.cache.Set(ctx, u)
	if err != nil {
		// 写入缓存失败，不返回错误，因为缓存不是必须的
		log.Println("write cache failed", err)
	}
	return u, nil
}

func (repo *UserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := repo.dao.QueryByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Phone:    u.Phone.String,
		Password: u.Password,
		Nickname: u.Nickname,
		Birthday: time.UnixMilli(u.Birthday),
		AboutMe:  u.AboutMe,
	}
}

func (repo *UserRepository) toEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
		Birthday: u.Birthday.UnixMilli(),
		AboutMe:  u.AboutMe,
		Nickname: u.Nickname,
	}
}
