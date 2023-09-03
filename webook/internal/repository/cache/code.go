package cache

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/allegro/bigcache"
	"github.com/redis/go-redis/v9"
)

const (
	tryGetLimits uint = 2
)

var (
	//go:embed lua/set_code.lua
	luaSetCode string
	//go:embed lua/verify_code.lua
	luaVerifyCode             string
	ErrCodeSendTooMany        = errors.New("发送验证码太频繁")
	ErrUnknownForCode         = errors.New("发送验证码遇到未知错误")
	ErrCodeVerifyTooManyTimes = errors.New("验证次数太多")
	ErrCodeInvalidated        = errors.New("验证码已失效")
)

type CodeCache interface {
	Set(ctx context.Context, biz string, phone string, code string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}

// RedisCodeCache 基于 Redis 的实现
type RedisCodeCache struct {
	redis redis.Cmdable
}

func NewRedisCodeCache(cmd redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		redis: cmd,
	}
}

// Set 如果该手机在该业务场景下，验证码不存在（都已经过期），那么发送
// 如果已经有一个验证码，但是发出去已经一分钟了，允许重发
// 如果已经有一个验证码，但是没有过期时间，说明有不知名错误
// 如果已经有一个验证码，但是发出去不到一分钟，不允许重发
// 验证码有效期 10 分钟
func (c *RedisCodeCache) Set(ctx context.Context, biz string, phone string, code string) error {
	res, err := c.redis.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		return nil
	case -1:
		//	最近发过
		return ErrCodeSendTooMany
	default:
		// 系统错误，比如说 -2，是 key 冲突
		// 其它响应码，不知道是啥鬼东西
		// TODO 按照道理，这里要考虑记录日志，但是我们暂时还没有日志模块，所以暂时不管
		return ErrUnknownForCode
	}
}

// Verify 验证验证码
// 如果验证码是一致的，那么删除
// 如果验证码不一致，那么保留的
func (c *RedisCodeCache) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	res, err := c.redis.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, inputCode).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		return true, nil
	case -1:
		//	验证次数耗尽，一般都是意味着有人在捣乱
		return false, ErrCodeVerifyTooManyTimes
	default:
		// 验证码不对
		return false, nil
	}
}

func (c *RedisCodeCache) key(biz string, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

// BigCodeCache 基于 BigCache 的实现
type BigCodeCache struct {
	bCache *bigcache.BigCache
}

type bCachedCode struct {
	Code               string `json:"code"`
	TryGets            uint   `json:"tryGets"`
	CreatedAtTimestamp int64  `json:"createdAtTimestamp"`
	Invalid            bool   `json:"invalid"`
}

func NewMemCodeCache(bCache *bigcache.BigCache) CodeCache {
	return &BigCodeCache{
		bCache: bCache,
	}
}

func (c *BigCodeCache) Set(ctx context.Context, biz string, phone string, code string) error {
	bs, err := c.bCache.Get(c.key(biz, phone))
	if err != nil {
		if !errors.Is(err, bigcache.ErrEntryNotFound) {
			return ErrUnknownForCode
		}
		// first time setting code
		ncc := bCachedCode{
			Code:               code,
			TryGets:            0,
			CreatedAtTimestamp: time.Now().UnixMilli(),
			Invalid:            false,
		}
		bs, err = json.Marshal(&ncc)
		if err != nil {
			return ErrUnknownForCode
		}
		err = c.bCache.Set(c.key(biz, phone), bs)
		if err != nil {
			return ErrUnknownForCode
		}
		return nil
	}

	// code already been set before
	var cc bCachedCode
	err = json.Unmarshal(bs, &cc)
	if err != nil {
		return ErrUnknownForCode
	}

	if cc.CreatedAtTimestamp == 0 {
		c.bCache.Delete(c.key(biz, phone))
		return ErrUnknownForCode
	}

	if (!cc.Invalid) && (time.Since(time.UnixMilli(cc.CreatedAtTimestamp)) < time.Second*60) {
		return ErrCodeSendTooMany
	}

	// code set before but has expired, need to reset the code
	ucc := bCachedCode{
		Code:               code,
		TryGets:            0,
		CreatedAtTimestamp: cc.CreatedAtTimestamp,
		Invalid:            false,
	}
	bs, err = json.Marshal(&ucc)
	if err != nil {
		return ErrUnknownForCode
	}
	err = c.bCache.Set(c.key(biz, phone), bs)
	return err
}

func (c *BigCodeCache) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {

	bs, err := c.bCache.Get(c.key(biz, phone))
	if err != nil {
		if !errors.Is(err, bigcache.ErrEntryNotFound) {
			// log it
			return false, ErrUnknownForCode
		}
		return false, nil
	}

	var cc bCachedCode
	err = json.Unmarshal(bs, &cc)
	if err != nil {
		return false, ErrUnknownForCode
	}

	if cc.Invalid {
		return false, ErrCodeInvalidated
	}

	if time.Since(time.UnixMilli(cc.CreatedAtTimestamp)) > time.Minute*10 {
		// code already expired
		cc.Invalid = true
		bs, err = json.Marshal(cc)
		if err != nil {
			return false, ErrUnknownForCode
		}
		err = c.bCache.Set(c.key(biz, phone), bs)
		if err != nil {
			return false, ErrUnknownForCode
		}
		return false, ErrCodeInvalidated
	}

	if cc.TryGets >= tryGetLimits {
		cc.Invalid = true
		bs, err = json.Marshal(cc)
		if err != nil {
			return false, ErrUnknownForCode
		}
		err = c.bCache.Set(c.key(biz, phone), bs)
		if err != nil {
			return false, ErrUnknownForCode
		}
		return false, ErrCodeVerifyTooManyTimes
	}

	if cc.Code != inputCode {
		// 验证码不对
		ucc := bCachedCode{
			Code:               cc.Code,
			TryGets:            cc.TryGets + 1,
			CreatedAtTimestamp: cc.CreatedAtTimestamp,
			Invalid:            false,
		}
		bs, err = json.Marshal(&ucc)
		if err != nil {
			return false, ErrUnknownForCode
		}
		err = c.bCache.Set(c.key(biz, phone), bs)
		if err != nil {
			return false, ErrUnknownForCode
		}
		return false, nil
	}

	// success, need to invalidate the used code
	cc.Invalid = true
	bs, err = json.Marshal(cc)
	if err != nil {
		return false, ErrUnknownForCode
	}
	err = c.bCache.Set(c.key(biz, phone), bs)
	if err != nil {
		return false, ErrUnknownForCode
	}

	return true, nil
}

func (c *BigCodeCache) key(biz string, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
