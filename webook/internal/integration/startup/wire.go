//go:build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/xiaoshanjiang/my-geektime/webook/internal/repository"
	article2 "github.com/xiaoshanjiang/my-geektime/webook/internal/repository/article"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/repository/cache"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/repository/dao"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/repository/dao/article"
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
var articlSvcProvider = wire.NewSet(
	article.NewGORMArticleDAO,
	article2.NewArticleRepository,
	service.NewArticleService,
)
var interactiveSvcProvider = wire.NewSet(
	service.NewInteractiveService,
	repository.NewCachedInteractiveRepository,
	dao.NewGORMInteractiveDAO,
	cache.NewRedisInteractiveCache,
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 最基础的第三方依赖
		thirdProvider,
		userSvcProvider,
		articlSvcProvider,
		// Cache 部分
		cache.NewRedisCodeCache,
		// repository 部分
		repository.NewCachedCodeRepository,
		// service 部分
		// 集成测试我们显式指定使用内存实现
		ioc.InitSmsMemoryService,

		// 指定啥也不干的 wechat service
		InitPhantomWechatService,
		service.NewSMSCodeService,
		// handler 部分
		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		web.NewArticleHandler,
		ijwt.NewRedisJWTHandler,

		// gin 的中间件
		ioc.InitMiddlewares,

		// Web 服务器
		ioc.InitWebServer,
	)
	// 随便返回一个
	return gin.Default()
}

func InitArticleHandler(dao article.ArticleDAO) *web.ArticleHandler {
	wire.Build(thirdProvider,
		//userSvcProvider,
		//cache.NewRedisArticleCache,
		//wire.InterfaceValue(new(article.ArticleDAO), dao),
		article2.NewArticleRepository,
		service.NewArticleService,
		web.NewArticleHandler)
	return new(web.ArticleHandler)
}

func InitUserSvc() service.UserService {
	wire.Build(thirdProvider, userSvcProvider)
	return service.NewUserService(nil, nil)
}

func InitJwtHdl() ijwt.Handler {
	wire.Build(thirdProvider, ijwt.NewRedisJWTHandler)
	return ijwt.NewRedisJWTHandler(nil)
}

func InitInteractiveService() service.InteractiveService {
	wire.Build(thirdProvider, interactiveSvcProvider)
	return service.NewInteractiveService(nil, nil)
}
