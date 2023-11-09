package ioc

import (
	"os"

	"github.com/xiaoshanjiang/my-geektime/webook/internal/service/oauth2/wechat"
	logger2 "github.com/xiaoshanjiang/my-geektime/webook/pkg/logger"
)

func InitWechatService(l logger2.LoggerV1) wechat.Service {
	os.Setenv("WECHAT_APP_ID", "wx7256bc69ab349c72") //TODO: remove this!
	os.Setenv("WECHAT_APP_SECRET", "secret")         //TODO: remove this!

	appId, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		panic("没有找到环境变量 WECHAT_APP_ID ")
	}
	appKey, ok := os.LookupEnv("WECHAT_APP_SECRET")
	if !ok {
		panic("没有找到环境变量 WECHAT_APP_SECRET")
	}
	// 692jdHsogrsYqxaUK9fgxw
	return wechat.NewService(appId, appKey, l)
}

// func NewWechatHandlerConfig() web.WechatHandlerConfig {
// 	return web.WechatHandlerConfig{
// 		Secure: false,
// 	}
// }
