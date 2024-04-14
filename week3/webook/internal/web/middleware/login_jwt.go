package middleware

import (
	"encoding/gob"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
	"webook/internal/web"
)

type LoginJWTMiddlewareBuild struct {
}

func (m *LoginJWTMiddlewareBuild) CheckLogin() gin.HandlerFunc {
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/user/signup" ||
			ctx.Request.URL.Path == "/user/login" {
			return
		}
		authCode := ctx.GetHeader("Authorization")
		if authCode == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
			return
		}
		authSeg := strings.Split(authCode, " ")
		if len(authSeg) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token格式错误"})
			return
		}
		tokenStr := authSeg[1]
		uc := web.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, uc, func(token *jwt.Token) (interface{}, error) {
			return web.JWTKey, nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token解析错误"})
			return
		}

		// 刷新token
		if token == nil || token.Valid {
			// token非法，或者过期
			ctx.AbortWithStatus(http.StatusUnauthorized)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token非法，或者过期"})
			return
		}

	}
}
