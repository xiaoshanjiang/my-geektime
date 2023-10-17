package service

import (
	"context"
	"errors"

	"github.com/xiaoshanjiang/my-geektime/webook/internal/domain"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/repository"
	"github.com/xiaoshanjiang/my-geektime/webook/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicate
var ErrInvalidUserOrPassword = errors.New("邮箱或者密码不正确")

type UserService interface {
	Login(ctx context.Context, email, password string) (domain.User, error)
	Signup(ctx context.Context, u domain.User) error
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	FindOrCreateByWechat(ctx context.Context, wechatInfo domain.WechatInfo) (domain.User, error)
	Profile(ctx context.Context, id int64) (domain.User, error)
	// UpdateNonSensitiveInfo 更新非敏感数据
	// 你可以在这里进一步补充究竟哪些数据会被更新
	UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error
}

type userService struct {
	repo repository.UserRepository
	l    logger.LoggerV1
}

func NewUserService(repo repository.UserRepository, l logger.LoggerV1) UserService {
	return &userService{
		repo: repo,
		l:    l,
	}
}

func NewUserServiceV1(repo repository.UserRepository, l *zap.Logger) UserService {
	return &userService{
		repo: repo,
		// 预留了变化空间
		//logger: zap.L(),
	}
}

func (svc *userService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *userService) Signup(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

// FindOrCreate 如果手机号不存在，那么会初始化一个用户
func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	// 这是一种优化写法, 大部分人会命中这个分支
	u, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		return u, err
	}
	// 要执行注册
	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	// 注册有问题，但是又不是用户手机号码冲突，说明是系统错误
	if err != nil && err != repository.ErrUserDuplicate {
		return domain.User{}, err
	}
	// 主从模式下，这里要从主库中读取，暂时我们不需要考虑
	return svc.repo.FindByPhone(ctx, phone)
}

func (svc *userService) FindOrCreateByWechat(ctx context.Context,
	info domain.WechatInfo) (domain.User, error) {
	u, err := svc.repo.FindByWechat(ctx, info.OpenID)
	if err != repository.ErrUserNotFound {
		return u, err
	}
	u = domain.User{
		WechatInfo: info,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil && err != repository.ErrUserDuplicate {
		return u, err
	}
	// 因为这里会遇到主从延迟的问题
	return svc.repo.FindByWechat(ctx, info.OpenID)
}

func (svc *userService) Profile(ctx context.Context, id int64) (domain.User, error) {
	return svc.repo.FindById(ctx, id)
}

func (svc *userService) UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error {
	user.Email = ""
	user.Phone = ""
	user.Password = ""
	return svc.repo.Update(ctx, user)
}
