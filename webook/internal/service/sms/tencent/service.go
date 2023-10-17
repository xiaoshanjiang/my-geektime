package tencent

import (
	"context"
	"fmt"

	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"github.com/xiaoshanjiang/my-geektime/webook/pkg/ratelimit"
	"go.uber.org/zap"
)

type Service struct {
	appId    *string
	signName *string
	client   *sms.Client
	limiter  ratelimit.Limiter
}

func NewService(c *sms.Client, appId string, signName string, limiter ratelimit.Limiter) *Service {
	return &Service{
		client:   c,
		appId:    ekit.ToPtr[string](appId),
		signName: ekit.ToPtr[string](signName),
		limiter:  limiter,
	}
}

func (s *Service) Send(ctx context.Context, biz string,
	args []string, numbers ...string) error {
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appId
	req.SignName = s.signName
	req.TemplateId = ekit.ToPtr[string](biz)
	req.PhoneNumberSet = toStringPtrSlice(numbers)
	req.TemplateParamSet = toStringPtrSlice(args)
	req.SetContext(ctx)
	resp, err := s.client.SendSms(req)
	zap.L().Debug("发送短信", zap.Any("req", req), zap.Any("resp", resp), zap.Error(err))
	if err != nil {
		return fmt.Errorf("腾讯短信服务发送失败 %w", err)
	}
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "Ok" {
			return fmt.Errorf("发送失败，code: %s, 原因：%s", *status.Code, *status.Message)
		}
	}
	return nil
}

func toStringPtrSlice(src []string) []*string {
	return slice.Map[string, *string](src, func(idx int, src string) *string {
		return &src
	})
}
