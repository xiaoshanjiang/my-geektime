package ioc

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/xiaoshanjiang/my-geektime/webook/internal/web"
	ijwt "github.com/xiaoshanjiang/my-geektime/webook/internal/web/jwt"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/web/middleware"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler,
	oauth2WechatHdl *web.OAuth2WechatHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	oauth2WechatHdl.RegisterRoutes(server)
	return server
}

func InitMiddlewares(redisClient redis.Cmdable, jwtHdl ijwt.Handler) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHandler(),
		middleware.NewLoginJWTMiddlewareBuilder(jwtHdl).Build(),
		// ratelimit.NewBuilder(redisClient, time.Second, 100).Build(),
	}
}

func corsHandler() gin.HandlerFunc {
	return cors.New(cors.Config{
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
	})
}

func sessionHandlerFunc() gin.HandlerFunc {
	//store := cookie.NewStore([]byte("secret"))

	// 这是基于内存的实现，第一个参数是 authentication key，最好是 32 或者 64 位
	// 第二个参数是 encryption key
	store := memstore.NewStore([]byte("moyn8y9abnd7q4zkq2m73yw8tu9j5ixm"),
		[]byte("o6jdlo2cb9f9pb6h46fjmllw481ldebj"))
	// 第一个参数是最大空闲连接数量
	// 第二个就是 tcp，你不太可能用 udp
	// 第三、四个 就是连接信息和密码
	// 第五第六就是两个 key
	//store, err := redis.NewStore(16, "tcp",
	//	"localhost:6379", "",
	//	// authentication key, encryption key
	//	[]byte("moyn8y9abnd7q4zkq2m73yw8tu9j5ixm"),
	//	[]byte("o6jdlo2cb9f9pb6h46fjmllw481ldebj"))
	//if err != nil {
	//	panic(err)
	//}

	// cookie 的名字叫做ssid
	return sessions.Sessions("ssid", store)
}
