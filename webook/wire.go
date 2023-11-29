//go:build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/xiaoshanjiang/my-geektime/webook/internal/events/article"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/repository"
	article2 "github.com/xiaoshanjiang/my-geektime/webook/internal/repository/article"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/repository/cache"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/repository/dao"
	article3 "github.com/xiaoshanjiang/my-geektime/webook/internal/repository/dao/article"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/service"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/web"
	ijwt "github.com/xiaoshanjiang/my-geektime/webook/internal/web/jwt"
	"github.com/xiaoshanjiang/my-geektime/webook/ioc"
)

func InitWebServer() *App {
	wire.Build(
		// 最基础的第三方依赖
		ioc.InitDB, ioc.InitRedis,
		ioc.InitLogger,
		ioc.InitKafka,
		ioc.NewConsumers,
		ioc.NewSyncProducer,

		// consumer
		article.NewInteractiveReadEventConsumer,
		article.NewKafkaProducer,

		// DAO 部分
		dao.NewGORMUserDAO,
		article3.NewGORMArticleDAO,
		dao.NewGORMInteractiveDAO,

		// Cache 部分
		cache.NewRedisInteractiveCache,
		cache.NewRedisUserCache,
		cache.NewRedisCodeCache,

		// repository 部分
		repository.NewCachedUserRepository,
		repository.NewCachedCodeRepository,
		repository.NewCachedInteractiveRepository,
		article2.NewArticleRepository,

		// service 部分
		// ioc.InitSmsService,
		// 直接基于内存实现
		ioc.InitSmsMemoryService,
		ioc.InitWechatService,
		service.NewUserService,
		service.NewSMSCodeService,
		service.NewArticleService,

		// handler 部分
		ijwt.NewRedisJWTHandler,
		web.NewUserHandler,
		web.NewArticleHandler,
		web.NewOAuth2WechatHandler,
		// ioc.NewWechatHandlerConfig,

		// gin 的中间件
		ioc.InitMiddlewares,

		// Web 服务器
		ioc.InitWebServer,
		// 组装我这个结构体的所有字段
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
