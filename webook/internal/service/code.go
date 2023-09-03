package service

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/xiaoshanjiang/my-geektime/webook/internal/repository"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/service/sms"
)

var (
	ErrCodeSendTooMany = repository.ErrCodeSendTooMany
	ErrCodeInvalidated = repository.ErrCodeInvalidated
)

const codeTplId = "1877556"

type CodeService interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}

// SMSCodeService 短信验证码的实现
type SMSCodeService struct {
	sms  sms.Service
	repo repository.CodeRepository
}

func NewSMSCodeService(svc sms.Service, repo repository.CodeRepository) CodeService {
	return &SMSCodeService{
		sms:  svc,
		repo: repo,
	}
}

// Send 生成一个随机验证码，并发送
func (c *SMSCodeService) Send(ctx context.Context, biz string, phone string) error {
	code := c.generateCode()
	err := c.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	err = c.sms.Send(ctx, codeTplId, []string{code}, phone)
	return err
}

// Verify 验证验证码
func (c *SMSCodeService) Verify(ctx context.Context,
	biz string,
	phone string,
	inputCode string) (bool, error) {
	ok, err := c.repo.Verify(ctx, biz, phone, inputCode)
	// 这里我们在 service 层面上对 Handler 屏蔽了最为特殊的错误
	if err == repository.ErrCodeVerifyTooManyTimes || err == repository.ErrCodeInvalidated {
		// 在接入了告警之后，这边要告警
		// 因为这意味着有人在搞你
		return false, nil
	}
	return ok, err
}

func (svc *SMSCodeService) generateCode() string {
	// 六位数，num 在 0, 999999 之间，包含 0 和 999999
	num := rand.Intn(1000000)
	// 不够六位的，加上前导 0
	// 000001
	return fmt.Sprintf("%06d", num)
}
