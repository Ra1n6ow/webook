package main

func main() {
	// db := initDB()
	// redisClient := initRedis()
	// server := initWebServer()
	// initUserHandler(server, db, redisClient)
	server := InitWebServer()
	server.Run(":8080")
}

// func initUserHandler(server *gin.Engine, db *gorm.DB, cmd redis.Cmdable) {
// 	userDAO := dao.NewUserDAO(db)
// 	userCache := cache.NewUserCache(cmd, time.Minute*15)
// 	codeCache := cache.NewCodeCache(cmd)
// 	userRepo := repository.NewUserRepository(userDAO, userCache)
// 	codeRepo := repository.NewCodeRepository(codeCache)
// 	userSvc := service.NewUserService(userRepo)
// 	localsms := localsms.NewService()
// 	codeSvc := service.NewCodeService(codeRepo, localsms)
// 	h := web.NewUserHandler(userSvc, codeSvc)
// 	h.RegisterRoutes(server)
// }

// func initWebServer() *gin.Engine {
// 	server := gin.Default()

// 	server.Use(cors.New(cors.Config{
// 		// AllowOrigins:     []string{"http://localhost:3000"},
// 		// AllowMethods:     []string{"POST", "GET"},
// 		// 允许前端发送给服务器的请求头
// 		AllowHeaders: []string{"content-type", "Authorization"},
// 		// 允许前端从服务器获取的响应头
// 		ExposeHeaders:    []string{"x-jwt-token"},
// 		AllowCredentials: true,
// 		AllowOriginFunc: func(origin string) bool {
// 			if strings.HasPrefix(origin, "http://localhost") {
// 				return true
// 			}
// 			return strings.Contains(origin, "company.com")
// 		},
// 		MaxAge: 1 * time.Hour,
// 	}))

// 	// 将 session 存储在 cookie 中， 一般是存储在 redis 中
// 	store := cookie.NewStore([]byte("secret12315"))
// 	server.Use(sessions.Sessions("ssid", store))
// 	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
// 		IgnorePaths("/users/login").
// 		IgnorePaths("/users/signup").
// 		IgnorePaths("/users/login_sms/code/send").
// 		IgnorePaths("/users/login_sms").
// 		Build())
// 	return server
// }
