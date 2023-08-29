package ioc

import (
	"github.com/redis/go-redis/v9"

	"github.com/xiaoshanjiang/my-geektime/webook/config"
)

func InitRedis() redis.Cmdable {
	rCfg := config.Config.Redis
	cmd := redis.NewClient(&redis.Options{
		Addr:     rCfg.Addr,
		Password: rCfg.Password,
		DB:       rCfg.DB,
	})
	return cmd
}
