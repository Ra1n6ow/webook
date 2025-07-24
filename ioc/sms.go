package ioc

import (
	"github.com/ra1n6ow/webook/internal/service/sms"
	"github.com/ra1n6ow/webook/internal/service/sms/localsms"
)

func InitSms() sms.Servicer {
	return localsms.NewService()
}
