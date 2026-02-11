package admin

import (
	"mall/adaptor"
	"mall/api"
	"mall/common"
	"mall/service/admin"

	"github.com/gin-gonic/gin"
)

type Ctrl struct {
	adaptor adaptor.IAdaptor
	user    *admin.Service
	hello   *admin.Service
}

func NewCtrl(adaptor adaptor.IAdaptor) *Ctrl {
	return &Ctrl{
		adaptor: adaptor,
		user:    admin.NewService(adaptor),
		hello:   admin.NewService(adaptor),
	}
}

func (ctrl *Ctrl) HelloWorld(ctx *gin.Context) {
	resp, errno := ctrl.hello.HelloWorld(ctx.Request.Context(), &common.AdminUser{}, nil)
	api.WriteResp(ctx, resp, errno)
}

func (c *Ctrl) GetUserInfo(ctx *gin.Context) {
	// token common.AdminUser
	resp, errno := c.user.GetUserInfo(ctx.Request.Context(), &common.AdminUser{})
	api.WriteResp(ctx, resp, errno)
}
