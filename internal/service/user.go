package service

import (
	"context"
	"errors"

	"github.com/ra1n6ow/webook/internal/domain"
	"github.com/ra1n6ow/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("用户名或密码错误")
	ErrRecordNotFound        = repository.ErrRecordNotFound
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
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
	if errors.Is(err, repository.ErrRecordNotFound) {
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
