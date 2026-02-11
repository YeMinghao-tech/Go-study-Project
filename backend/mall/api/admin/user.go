package admin

import (
	"mall/api"
	"mall/common"

	"github.com/gin-gonic/gin"
	//"mall/service/dto"
)

func (c *Ctrl) GetUserInfo(ctx *gin.Context) {
	user := api.GetAdminUserFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	resp, errno := c.user.GetUserInfo(ctx.Request.Context(), &common.AdminUser{})
	api.WriteResp(ctx, resp, errno)
}
