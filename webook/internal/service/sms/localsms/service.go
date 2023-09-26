package localsms

import (
	"context"
	"log"

	"github.com/xiaoshanjiang/my-geektime/webook/pkg/ratelimit"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	log.Println("验证码是", args)
	return nil
}

func (s *Service) GetLimiter() ratelimit.Limiter {
	return nil
}
