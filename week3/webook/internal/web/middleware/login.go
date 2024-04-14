package middleware

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuild struct {
}

func (m *LoginMiddlewareBuild) CheckLogin() gin.HandlerFunc {
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/user/signup" ||
			ctx.Request.URL.Path == "/user/login" {
			return
		}
		sess := sessions.Default(ctx)
		uid := sess.Get("userId")
		if uid == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
			return
		}
		ctx.Set("userId", uid)

		// 刷新ssid
		now := time.Now()

		const updateTimeKey = "update_time"
		t := sess.Get(updateTimeKey)
		lastTime, ok := t.(time.Time)
		if t == nil || !ok || lastTime.Sub(time.Now()) > time.Minute {
			// 刷新ssid
			sess.Set(updateTimeKey, now)
			sess.Set("userId", uid)
			err := sess.Save()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
