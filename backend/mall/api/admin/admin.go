package admin

import (
	"mall/adaptor"
	"mall/api"
	"mall/common"
	"mall/service/admin"

	"github.com/gin-gonic/gin"
)

type Ctrl struct {
	adaptor *adaptor.Adaptor
	user    *admin.Service
}

func NewCtrl(adaptor *adaptor.Adaptor) *Ctrl {
	return &Ctrl{
		adaptor: adaptor,
		user:    admin.NewService(adaptor),
	}
}

func (ctrl *Ctrl) HelloWorld(ctx *gin.Context) {
	api.WriteResp(ctx, "helloworld", common.OK)
}
