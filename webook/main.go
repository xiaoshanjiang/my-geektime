package main

import (
	"GeekTime/my-geektime/webook/internal/web"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	// v1 := server.Group("/v1")
	// users := server.Group("/users/v1")
	// u.RegisterRoutesV1(server.Group(("/users")))
	u := web.NewUserHandler()
	u.RegisterRoutes(server)
	server.Run(":8080")
}
