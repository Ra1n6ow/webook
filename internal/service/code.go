package service

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/ra1n6ow/webook/internal/repository"
	"github.com/ra1n6ow/webook/internal/service/sms"
)

var ErrCodeSendTooMany = repository.ErrCodeVerifyTooMany

type CodeServicer interface {
	Send(ctx context.Context, biz, phone string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type CodeService struct {
	repo repository.CodeRepositoryer
	sms  sms.Servicer
}

func NewCodeService(repo repository.CodeRepositoryer, sms sms.Servicer) CodeServicer {
	return &CodeService{
		repo: repo,
		sms:  sms,
	}
}

func (svc *CodeService) Send(ctx context.Context, biz, phone string) error {
	code := svc.generate()
	err := svc.repo.Set(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	// 发送验证码的模板 id
	const codeTplId = "1877556"
	return svc.sms.Send(ctx, codeTplId, []string{code}, phone)
}

func (svc *CodeService) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	ok, err := svc.repo.Verify(ctx, biz, phone, inputCode)
	if err == repository.ErrCodeVerifyTooMany {
		// 相当于，我们对外面屏蔽了验证次数过多的错误，我们就是告诉调用者，你这个不对
		return false, nil
	}
	return ok, err
}

func (svc *CodeService) generate() string {
	// 0-999999
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}
