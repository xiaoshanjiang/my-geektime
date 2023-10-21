package startup

import (
	"github.com/xiaoshanjiang/my-geektime/webook/pkg/logger"
)

func InitLog() logger.LoggerV1 {
	return logger.NewNoOpLogger()
}
