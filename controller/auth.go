package controller

import (
	"encoding/json"
	"time"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jeffcail/gateway-web/dao"
	"github.com/jeffcail/gateway-web/dto"
	"github.com/jeffcail/gateway-web/middleware"
	"github.com/jeffcail/gateway-web/public"
)

type AdminLoginHandler struct{}

func RegisterAdmin(group *gin.RouterGroup) {
	a := &AdminLoginHandler{}
	group.POST("/login", a.AdminLogin)
	group.GET("/logout", a.AdminLogout)
}

// AdminLogin godoc
// @Summary 管理员登陆
// @Description 管理员登陆
// @Tags 管理员接口
// @ID /admin_login/login
// @Accept json
// @Produce json
// @Param body body dto.AdminLoginInput true "body"
// @Success 200 {object} middleware.Response{data=dto.AdminLoginOutput} "success"
// @Router /admin_login/login [post]
func (admin *AdminLoginHandler) AdminLogin(c *gin.Context) {
	params := &dto.AdminLoginInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 1001, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	ad := &dao.GatewayAdmin{}
	ad, err = ad.CheckAdminLogin(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	sessInfo := &dto.AdminSessionInfo{
		ID:        ad.ID,
		UserName:  ad.UserName,
		LoginTime: time.Now(),
	}
	sBits, _ := json.Marshal(sessInfo)

	sess := sessions.Default(c)
	sess.Set(public.AdminSessionInfoKey, string(sBits))
	sess.Save()

	out := &dto.AdminLoginOutput{Token: ad.UserName}
	middleware.ResponseSuccess(c, out)
}

// AdminLogout godoc
// @Summary 管理员退出
// @Description 管理员退出
// @Tags 管理员接口
// @ID /admin_login/logout
// @Accept json
// @Produce json
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin_login/logout [get]
func (admin *AdminLoginHandler) AdminLogout(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Delete(public.AdminSessionInfoKey)
	sess.Save()
	middleware.ResponseSuccess(c, "")
}
