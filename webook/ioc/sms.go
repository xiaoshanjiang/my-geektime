package ioc

import (
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentSMS "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"

	"github.com/xiaoshanjiang/my-geektime/webook/internal/service/sms"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/service/sms/localsms"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/service/sms/tencent"
	"github.com/xiaoshanjiang/my-geektime/webook/pkg/ratelimit"
)

func InitSmsService(cmd redis.Cmdable) sms.Service {
	//return initSmsTencentService()
	return initSmsTencentService()
}

func initSmsTencentService() sms.Service {
	secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	if !ok {
		panic("没有找到环境变量 SMS_SECRET_ID ")
	}
	secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")

	c, err := tencentSMS.NewClient(common.NewCredential(secretId, secretKey),
		"ap-nanjing",
		profile.NewClientProfile())
	if !ok || err != nil {
		panic("没有找到环境变量 SMS_SECRET_KEY")
	}
	redis := InitRedis()
	ratelimiter := ratelimit.NewRedisSlidingWindowLimiter(redis, time.Minute, 1000)
	return tencent.NewService(c, "1400842696", "妙影科技", ratelimiter)
}

// InitSmsMemoryService 使用基于内存，输出到控制台的实现
func InitSmsMemoryService() sms.Service {
	return localsms.NewService()
}
