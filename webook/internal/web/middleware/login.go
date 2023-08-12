package middleware

import (
	"encoding/gob"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// LoginMiddlewareBuilder 扩展性
type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	// 用 Go 的方式编码解码
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		// 不需要登录校验的
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		if id == nil {
			// 没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		sess.Set("userId", id)
		sess.Options(sessions.Options{
			MaxAge: 60,
		})
		updateTime := sess.Get("update_time")
		now := time.Now()
		// 说明还没有刷新过，刚登陆，还没刷新过
		if updateTime == nil {
			sess.Set("update_time", now)
			if err := sess.Save(); err != nil {
				panic(err)
			}
			return
		}
		// updateTime 是有的
		updateTimeVal, _ := updateTime.(time.Time)
		if now.Sub(updateTimeVal) > time.Second*10 {
			sess.Set("update_time", now)
			if err := sess.Save(); err != nil {
				panic(err)
			}
		}
	}
}
