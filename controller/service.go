package controller

import (
	"errors"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/jeffcail/gateway-web/dao"
	"github.com/jeffcail/gateway-web/dto"
	"github.com/jeffcail/gateway-web/middleware"
	"github.com/jeffcail/gateway-web/public"
	"strings"
)

type ServiceController struct{}

func ServiceRegister(group *gin.RouterGroup) {
	service := &ServiceController{}

	group.POST("/service_add_tcp", service.ServiceAddTcp)
}

// ServiceAddTcp godoc
// @Summary tcp服务添加
// @Description tcp服务添加
// @Tags 服务管理
// @ID /service/service_add_tcp
// @Accept json
// @Produce json
// @Param body body dto.ServiceAddTcpInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_add_tcp [post]
func (service *ServiceController) ServiceAddTcp(c *gin.Context) {
	param := &dto.ServiceAddTcpInput{}
	if err := param.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 20001, err)
		return
	}

	// 验证 service_name 是否被占用
	infoSearch := &dao.GatewayServiceInfo{
		ServiceName: param.ServiceName,
		IsDelete:    0,
	}
	_, err := infoSearch.Find(c, lib.GORMDefaultPool, infoSearch)
	if err != nil {
		middleware.ResponseError(c, 2002, errors.New("服务名被占用，重新输入"))
		return
	}

	// 验证端口是否被占用
	tcpRuleSearch := &dao.TcpRule{
		Port: param.Port,
	}
	if _, err := tcpRuleSearch.Find(c, lib.GORMDefaultPool, tcpRuleSearch); err != nil {
		middleware.ResponseError(c, 2003, errors.New("服务端口被占用,请重新输入"))
		return
	}
	grpcRuleSearch := &dao.GrpcRule{Port: param.Port}
	if _, err := grpcRuleSearch.Find(c, lib.GORMDefaultPool, grpcRuleSearch); err != nil {
		middleware.ResponseError(c, 2004, errors.New("服务端口被占用,请重新输入"))
		return
	}

	// 验证ip 与 权重数量的一致性
	if len(strings.Split(param.IpList, "")) != len(strings.Split(param.WeightList, "")) {
		middleware.ResponseError(c, 2005, errors.New("ip列表与权重列表设置不匹配"))
		return
	}

	tx := lib.GORMDefaultPool.Begin()
	info := &dao.GatewayServiceInfo{LoadType: public.LoadTypeTcp, ServiceName: param.ServiceName, ServiceDesc: param.ServiceDesc}
	err = info.Save(c, tx)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}

	// 轮询信息保存
	loadBalance := &dao.GatewayServiceLoadBalance{
		ServiceID:  info.ID,
		RoundType:  param.RoundType,
		IPList:     param.IpList,
		WeightList: param.WeightList,
		ForbidList: param.ForbidList,
	}
	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}

	// 保存Tcp rule 信息
	tcpRule := &dao.TcpRule{ServiceID: info.ID, Port: param.Port}
	if err := tcpRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2008, err)
		return
	}

	// 网关权限信息
	accessControl := &dao.GatewayServiceAccessControl{
		ServiceID:         info.ID,
		OpenAuth:          param.OpenAuth,
		BlackList:         param.BlackList,
		WhiteList:         param.WhiteList,
		WhiteHostName:     param.WhiteHostName,
		ClientipFlowLimit: param.ClientIPFlowLimit,
		ServiceFlowLimit:  param.ServiceFlowLimit,
	}
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2009, err)
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(c, "")
	return
}
