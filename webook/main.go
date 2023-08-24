package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"GeekTime/my-geektime/webook/config"
	"GeekTime/my-geektime/webook/internal/repository"
	"GeekTime/my-geektime/webook/internal/repository/dao"
	"GeekTime/my-geektime/webook/internal/service"
	"GeekTime/my-geektime/webook/internal/web"
	"GeekTime/my-geektime/webook/internal/web/middleware"
)

func main() {
	db := initDB()
	server := initWebServer()

	u := initUser(db)
	u.RegisterRoutes(server)

	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello! Welcome to webook!")
	})
	server.Run(":8080")
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	// redisClient := redis.NewClient(&redis.Options{
	// 	Addr: config.Config.Redis.Addr,
	// })
	// server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())

	server.Use(cors.New(cors.Config{
		//AllowOrigins: []string{"*"},
		//AllowMethods: []string{"POST", "GET"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// 你不加这个，前端是拿不到的
		ExposeHeaders: []string{"x-jwt-token"},
		// 是否允许你带 cookie 之类的东西
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				// 你的开发环境
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	// 步骤1
	//store := cookie.NewStore([]byte("secret"))

	// store := memstore.NewStore(
	// 	[]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"),
	// 	[]byte("0Pf2r0wZBpXVXlQNdpwCXN4ncnlnZSc3"),
	// )
	// store, err := redis.NewStore(
	// 	16,
	// 	"tcp", "localhost:6379", "",
	// 	[]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"),
	// 	[]byte("0Pf2r0wZBpXVXlQNdpwCXN4ncnlnZSc3"),
	// )

	// if err != nil {
	// 	panic(err)
	// }

	//myStore := &sqlx_store.Store{}

	// server.Use(sessions.Sessions("mysession", store))
	// 步骤3
	//server.Use(middleware.NewLoginMiddlewareBuilder().
	//	IgnorePaths("/users/signup").
	//	IgnorePaths("/users/login").Build())
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/hello").
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login").Build())

	// v1
	//middleware.IgnorePaths = []string{"sss"}
	//server.Use(middleware.CheckLogin())

	// 不能忽略sss这条路径
	//server1 := gin.Default()
	//server1.Use(middleware.CheckLogin())
	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                      config.Config.DB.DSN,
		DefaultStringSize:        256,  // default size for string fields
		DisableDatetimePrecision: true, // disable datetime precision, which not supported before MySQL 5.6
	}), &gorm.Config{
		CreateBatchSize: 1000,
	})

	if err != nil {
		// 我只会在初始化过程中 panic
		// panic 相当于整个 goroutine 结束
		// 一旦初始化过程出错，应用就不要启动了
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
