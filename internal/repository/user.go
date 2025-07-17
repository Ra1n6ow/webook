package repository

import (
	"context"
	"log"

	"github.com/ra1n6ow/webook/internal/domain"
	"github.com/ra1n6ow/webook/internal/repository/cache"
	"github.com/ra1n6ow/webook/internal/repository/dao"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrRecordNotFound = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, cache *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: cache,
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	ud := dao.User{
		Email:    u.Email,
		Password: u.Password,
	}
	return r.dao.Insert(ctx, ud)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.QueryByEmail(ctx, email)
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

	return r.dao.Update(ctx, ud)
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
