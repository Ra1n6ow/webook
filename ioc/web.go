package ioc

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/ra1n6ow/webook/internal/web"
	"github.com/ra1n6ow/webook/internal/web/middleware"
)

func InitWebServer(m []gin.HandlerFunc, h *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(m...)
	h.RegisterRoutes(server)
	return server
}

func InitMiddlewares(InitSessionMiddleware gin.HandlerFunc) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			// AllowOrigins:     []string{"http://localhost:3000"},
			// AllowMethods:     []string{"POST", "GET"},
			// 允许前端发送给服务器的请求头
			AllowHeaders: []string{"content-type", "Authorization"},
			// 允许前端从服务器获取的响应头
			ExposeHeaders:    []string{"x-jwt-token"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				if strings.HasPrefix(origin, "http://localhost") {
					return true
				}
				return strings.Contains(origin, "company.com")
			},
			MaxAge: 1 * time.Hour,
		}),
		middleware.NewLoginJWTMiddlewareBuilder().
			IgnorePaths("/users/login").
			IgnorePaths("/users/signup").
			IgnorePaths("/users/login_sms/code/send").
			IgnorePaths("/users/login_sms").
			Build(),
		InitSessionMiddleware,
	}
}

func InitSessionMiddleware() gin.HandlerFunc {
	store := cookie.NewStore([]byte("secret12315"))
	return sessions.Sessions("ssid", store)
}
