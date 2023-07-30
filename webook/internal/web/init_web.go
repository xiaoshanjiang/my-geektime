package web

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes() *gin.Engine {
	server := gin.Default()
	registerUsersRoutes(server)
	return server
}

func registerUsersRoutes(server *gin.Engine) {
	u := &UserHandler{}
	server.POST("/users/signup", u.SignUp)
	server.POST("/users/login", u.Login)
	server.POST("/users/edit", u.Edit)
	server.GET("/users/profile", u.Profile)

	// If going RESTful style
	// server.POST("/user", u.SignUp)   // create user / sign up
	// server.GET("/users/:id", u.Profile)   // read
	// server.PUT("/users/:id", u.Edit) // ddit
	// server.POST("/users/login", u.Login)   // login
}
