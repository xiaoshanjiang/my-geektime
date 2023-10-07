//go:build wireinject

package integration

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/xiaoshanjiang/my-geektime/webook/internal/repository"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/repository/cache"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/repository/dao"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/service"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/web"
	ijwt "github.com/xiaoshanjiang/my-geektime/webook/internal/web/jwt"
	"github.com/xiaoshanjiang/my-geektime/webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 最基础的第三方依赖
		ioc.InitDB, ioc.InitRedis,

		// DAO 部分
		dao.NewGORMUserDAO,

		// Cache 部分
		cache.NewRedisUserCache, cache.NewRedisCodeCache,

		// repository 部分
		repository.NewCachedUserRepository,
		repository.NewCachedCodeRepository,

		// service 部分
		// ioc.InitSmsService,
		service.NewUserService,
		service.NewSMSCodeService,
		// 直接基于内存实现
		ioc.InitSmsMemoryService,
		ioc.InitWechatService,

		// handler 部分
		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		ioc.NewWechatHandlerConfig,
		ijwt.NewRedisJWTHandler,

		// gin 的中间件
		ioc.InitMiddlewares,

		// Web 服务器
		ioc.InitWebServer,
	)
	// 随便返回一个
	return gin.Default()
}
