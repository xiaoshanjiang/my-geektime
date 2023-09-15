//go:build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/xiaoshanjiang/my-geektime/webook/internal/repository"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/repository/cache"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/repository/dao"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/service"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/web"
	"github.com/xiaoshanjiang/my-geektime/webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitRedis, ioc.InitDB,

		// DAO 部分
		dao.NewGORMUserDAO,

		// Cache 部分
		cache.NewRedisUserCache, cache.NewRedisCodeCache,

		// repository 部分
		repository.NewCachedUserRepository,
		repository.NewCachedCodeRepository,

		// service 部分
		ioc.InitSmsService,
		service.NewSMSCodeService,
		service.NewUserService,

		// handler 部分
		web.NewUserHandler,

		// gin 的中间件
		ioc.InitMiddlewares,

		// Web 服务器
		ioc.InitWebServer,
	)
	// 随便返回一个
	return gin.Default()
}
