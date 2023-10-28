//go:build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/xiaoshanjiang/my-geektime/webook/internal/repository"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/repository/article"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/repository/cache"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/repository/dao"
	article2 "github.com/xiaoshanjiang/my-geektime/webook/internal/repository/dao/article"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/service"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/web"
	ijwt "github.com/xiaoshanjiang/my-geektime/webook/internal/web/jwt"
	"github.com/xiaoshanjiang/my-geektime/webook/ioc"
)

var thirdProvider = wire.NewSet(InitRedis, InitTestDB, InitLog)
var userSvcProvider = wire.NewSet(
	dao.NewGORMUserDAO,
	cache.NewRedisUserCache,
	repository.NewCachedUserRepository,
	service.NewUserService,
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 最基础的第三方依赖
		thirdProvider,
		userSvcProvider,
		// Cache 部分
		cache.NewRedisCodeCache,

		article2.NewGORMArticleDAO,
		// repository 部分
		repository.NewCachedCodeRepository,
		article.NewArticleRepository,
		// service 部分
		// 集成测试我们显式指定使用内存实现
		ioc.InitSmsMemoryService,

		// 指定啥也不干的 wechat service
		InitPhantomWechatService,
		service.NewSMSCodeService,
		service.NewArticleService,
		// handler 部分
		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		web.NewArticleHandler,
		InitWechatHandlerConfig,
		ijwt.NewRedisJWTHandler,

		// gin 的中间件
		ioc.InitMiddlewares,

		// Web 服务器
		ioc.InitWebServer,
	)
	// 随便返回一个
	return gin.Default()
}

func InitArticleHandler() *web.ArticleHandler {
	wire.Build(thirdProvider,
		article2.NewGORMArticleDAO,
		service.NewArticleService,
		web.NewArticleHandler,
		article.NewArticleRepository,
	)
	return &web.ArticleHandler{}
}

func InitUserSvc() service.UserService {
	wire.Build(thirdProvider, userSvcProvider)
	return service.NewUserService(nil, nil)
}

func InitJwtHdl() ijwt.Handler {
	wire.Build(thirdProvider, ijwt.NewRedisJWTHandler)
	return ijwt.NewRedisJWTHandler(nil)
}
