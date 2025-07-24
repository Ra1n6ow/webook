package service

import (
	"context"
	"errors"

	"github.com/ra1n6ow/webook/internal/domain"
	"github.com/ra1n6ow/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicateUser
	ErrInvalidUserOrPassword = errors.New("用户名或密码错误")
	ErrRecordNotFound        = repository.ErrUserNotFound
)

type UserServicer interface {
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email, password string) (domain.User, error)
	Edit(ctx context.Context, u domain.User) error
	Profile(ctx context.Context, userId int64) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
}

type UserService struct {
	repo repository.UserRepositoryer
}

func NewUserService(repo repository.UserRepositoryer) UserServicer {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) SignUp(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return s.repo.Create(ctx, u)
}

func (s *UserService) Login(ctx context.Context, email, password string) (domain.User, error) {
	u, err := s.repo.FindByEmail(ctx, email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (s *UserService) Edit(ctx context.Context, u domain.User) error {
	return s.repo.Update(ctx, u)
}

func (s *UserService) Profile(ctx context.Context, userId int64) (domain.User, error) {
	return s.repo.FindById(ctx, userId)
}

func (svc *UserService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	// 先找一下，我们认为，大部分用户是已经存在的用户
	u, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		// 有两种情况
		// err == nil, u 是可用的
		// err != nil，系统错误，
		return u, err
	}
	// 用户没找到
	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	// 有两种可能，一种是 err 恰好是唯一索引冲突（phone）
	// 一种是 err != nil，系统错误
	if err != nil && err != repository.ErrDuplicateUser {
		return domain.User{}, err
	}
	// 要么 err ==nil，要么ErrDuplicateUser，也代表用户存在
	// 主从延迟，理论上来讲，强制走主库
	return svc.repo.FindByPhone(ctx, phone)
}
