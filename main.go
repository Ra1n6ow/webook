package main

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/ra1n6ow/webook/internal/repository"
	"github.com/ra1n6ow/webook/internal/repository/dao"
	"github.com/ra1n6ow/webook/internal/service"
	"github.com/ra1n6ow/webook/internal/web"
	"github.com/ra1n6ow/webook/internal/web/middleware"
)

func main() {
	db := initDB()
	server := initWebServer()
	initUserHandler(server, db)
	server.Run(":8080")
}

func initUserHandler(server *gin.Engine, db *gorm.DB) {
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	h := web.NewUserHandler(svc)
	h.RegisterRoutes(server)
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:123456@tcp(localhost:3306)/webook?parseTime=true"))
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
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
	}))

	// 将 session 存储在 cookie 中
	store := cookie.NewStore([]byte("secret12315"))
	server.Use(sessions.Sessions("ssid", store))
	server.Use(middleware.NewLoginMiddlewareBuilder().IgnorePaths("/users/login").IgnorePaths("/users/signup").Build())
	return server
}
