package dto

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jeffcail/gateway-web/public"
)

type AdminInfoOutput struct {
	ID           int64     `json:"id"`
	UserName     string    `json:"user_name"`
	LoginTime    time.Time `json:"login_time"`
	Avatar       string    `json:"avatar"`
	Introduction string    `json:"introduction"`
	Roles        []string  `json:"roles"`
}

type ChangePwdInput struct {
	Password string `json:"password" form:"password" comment:"密码" example:"123456" validate:"required"`
}

func (param *ChangePwdInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}
