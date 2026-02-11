package router

import (
	"context"
	"mall/common"
	"mall/consts"
	"net/http"

	//"mall/consts"
	//"net/http"

	"github.com/gin-gonic/gin"
)

type TokenFun func(ctx context.Context, token string) (*common.User, error)
type TokenAdminFun func(ctx context.Context, token string) (*common.AdminUser, error)

// 用户侧鉴权中间件
func AuthMiddleware(filter func(*gin.Context) bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if filter != nil && !filter(c) {
			c.Next()
			return
		}
		// 这里是鉴权中间件
		c.Next()
	}
}

// 管理后台用户中间件
func AdminAuthMiddleware(filter func(*gin.Context) bool, getTokenFun TokenAdminFun) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if filter != nil && !filter(ctx) {
			ctx.Next()
			return
		}
		token := ctx.GetHeader(consts.AdminTokenKey)
		if len(token) == 0 {
			ctx.JSON(http.StatusUnauthorized, common.AuthErr)
			ctx.Abort()
			return
		}
		user, err := getTokenFun(ctx, token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, common.AuthErr.WithErr(err))
			ctx.Abort()
			return
		}
		ctx.Set(consts.AdminUserKey, user)
		ctx.Next()
	}
}
