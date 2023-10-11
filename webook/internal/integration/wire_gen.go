// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package integration

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/repository"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/repository/cache"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/repository/dao"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/service"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/web"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/web/jwt"
	"github.com/xiaoshanjiang/my-geektime/webook/ioc"
)

// Injectors from wire.go:

func InitWebServer() *gin.Engine {
	cmdable := ioc.InitRedis()
	handler := jwt.NewRedisJWTHandler(cmdable)
	v := ioc.InitMiddlewares(cmdable, handler)
	db := ioc.InitDB()
	userDAO := dao.NewGORMUserDAO(db)
	userCache := cache.NewRedisUserCache(cmdable)
	userRepository := repository.NewCachedUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository)
	smsService := ioc.InitSmsMemoryService()
	codeCache := cache.NewRedisCodeCache(cmdable)
	codeRepository := repository.NewCachedCodeRepository(codeCache)
	codeService := service.NewSMSCodeService(smsService, codeRepository)
	userHandler := web.NewUserHandler(userService, codeService, handler)
	wechatService := ioc.InitWechatService()
	wechatHandlerConfig := ioc.NewWechatHandlerConfig()
	oAuth2WechatHandler := web.NewOAuth2WechatHandler(wechatService, userService, handler, wechatHandlerConfig)
	engine := ioc.InitWebServer(v, userHandler, oAuth2WechatHandler)
	return engine
}