//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/ra1n6ow/webook/internal/repository"
	"github.com/ra1n6ow/webook/internal/repository/cache"
	"github.com/ra1n6ow/webook/internal/repository/dao"
	"github.com/ra1n6ow/webook/internal/service"
	"github.com/ra1n6ow/webook/internal/web"
	"github.com/ra1n6ow/webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB, ioc.InitRedis,

		dao.NewUserDAO, cache.NewCodeCache, cache.NewUserCache,

		repository.NewUserRepository, repository.NewCodeRepository,

		service.NewUserService, service.NewCodeService, ioc.InitSms,

		web.NewUserHandler,

		ioc.InitSessionMiddleware, ioc.InitMiddlewares, ioc.InitWebServer,
	)
	return new(gin.Engine)
}
