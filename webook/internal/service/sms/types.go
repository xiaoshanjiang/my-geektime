package sms

import "context"

// Service 发送短信的抽象
// 目前你可以理解为，这是一个为了适配不同的短信供应商的抽象
type Service interface {
	Send(ctx context.Context, biz string, args []string, numbers ...string) error
	//SendV1(ctx context.Context, tpl string, args []NamedArg, numbers ...string) error
	// 调用者需要知道实现者需要什么类型的参数，是 []string，还是 map[string]string
	//SendV2(ctx context.Context, tpl string, args any, numbers ...string) error
	//SendV3(ctx context.Context, tpl string, args T, numbers ...string) error
}

type NamedArg struct {
	Val  string
	Name string
}
