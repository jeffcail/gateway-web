package controller

import (
	"errors"
	"fmt"
	"strings"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/jeffcail/gateway-web/dao"
	"github.com/jeffcail/gateway-web/dto"
	"github.com/jeffcail/gateway-web/middleware"
	"github.com/jeffcail/gateway-web/public"
)

type ServiceController struct{}

func ServiceRegister(group *gin.RouterGroup) {
	service := &ServiceController{}
	group.GET("/service_list", service.ServiceList)
	group.GET("/service_detail", service.ServiceDetail)
	group.GET("/service_delete", service.ServiceDelete)

	group.POST("/service_add_tcp", service.ServiceAddTcp)
}

// ServiceDelete godoc
// @Summary 服务删除
// @Description 服务删除
// @Tags 服务管理
// @ID /service/service_delete
// @Accept json
// @Produce json
// @Param id query string true "服务ID"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_delete [get]
func (service *ServiceController) ServiceDelete(c *gin.Context) {
	params := &dto.ServiceDDInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	serviceInfo := &dao.GatewayServiceInfo{ID: params.ID}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	serviceInfo.IsDelete = 1
	if err := serviceInfo.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	middleware.ResponseSuccess(c, "")
}

// ServiceDetail godoc
// @Summary 服务详情
// @Description 服务详情
// @Tags 服务管理
// @ID /service/service_detail
// @Accept  json
// @Produce  json
// @Param id query string true "服务ID"
// @Success 200 {object} middleware.Response{data=dao.ServiceDetail} "success"
// @Router /service/service_detail [get]
func (service *ServiceController) ServiceDetail(c *gin.Context) {
	params := &dto.ServiceDDInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	serviceInfo := &dao.GatewayServiceInfo{ID: params.ID}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	serviceDetail, err := serviceInfo.ServiceDetail(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	middleware.ResponseSuccess(c, serviceDetail)
}

// ServiceList godoc
// @Summary 服务列表
// @Description 服务列表
// @Tags 服务管理
// @ID /service/service_list
// @Accept  json
// @Produce  json
// @Param info query string false "关键词"
// @Param page_size query int true "每页个数"
// @Param page_no query int true "当前页数"
// @Success 200 {object} middleware.Response{data=dto.ServiceListOutput} "success"
// @Router /service/service_list [get]
func (service *ServiceController) ServiceList(c *gin.Context) {
	params := &dto.ServiceListInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	serviceInfo := &dao.GatewayServiceInfo{}
	list, total, err := serviceInfo.PageList(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	outList := []dto.ServiceListItemOutput{}
	for _, item := range list {
		serviceDetail, err := item.ServiceDetail(c, tx, &item)
		if err != nil {
			middleware.ResponseError(c, 2003, err)
			return
		}
		//1、http后缀接入 clusterIP+clusterPort+path
		//2、http域名接入 domain
		//3、tcp、grpc接入 clusterIP+servicePort
		serviceAddr := "unknow"
		clusterIP := lib.GetStringConf("base.cluster.cluster_ip")
		clusterPort := lib.GetStringConf("base.cluster.cluster_port")
		clusterSSLPort := lib.GetStringConf("base.cluster.cluster_ssl_port")
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL &&
			serviceDetail.HTTPRule.NeedHttps == 1 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterSSLPort, serviceDetail.HTTPRule.Rule)
		}
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL &&
			serviceDetail.HTTPRule.NeedHttps == 0 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterPort, serviceDetail.HTTPRule.Rule)
		}
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypeDomain {
			serviceAddr = serviceDetail.HTTPRule.Rule
		}
		if serviceDetail.Info.LoadType == public.LoadTypeTcp {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIP, serviceDetail.TcpRule.Port)
		}
		if serviceDetail.Info.LoadType == public.LoadTypeGRPC {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIP, serviceDetail.GrpcRule.Port)
		}
		ipList := serviceDetail.LoadBalance.GetIpListByModel()
		counter, err := public.FlowCounterHandler.GetCounter(public.FlowServicePrefix + item.ServiceName)
		if err != nil {
			middleware.ResponseError(c, 2004, err)
			return
		}
		outItem := dto.ServiceListItemOutput{
			ID:          item.ID,
			ServiceName: item.ServiceName,
			ServiceDesc: item.ServiceDesc,
			LoadType:    item.LoadType,
			ServiceAddr: serviceAddr,
			Qps:         counter.Qps,
			Qpd:         counter.TotalCount,
			TotalNode:   len(ipList),
		}
		outList = append(outList, outItem)
	}
	out := &dto.ServiceListOutput{
		Total: total,
		List:  outList,
	}
	middleware.ResponseSuccess(c, out)
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
