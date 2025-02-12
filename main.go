package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ra1n6ow/webook/internal/web"
)

func main() {
	server := gin.Default()
	u := web.NewUserHandler()
	u.RegisterRoutes(server)
	server.Run(":8080")
}
