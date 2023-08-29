package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	server := InitWebServer()
	// 注册路由
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello, world")
	})
	server.Run(":8080")
}
