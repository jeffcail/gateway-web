package controller

import (
	"errors"
	"time"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/jeffcail/gateway-web/dao"
	"github.com/jeffcail/gateway-web/dto"
	"github.com/jeffcail/gateway-web/middleware"
	"github.com/jeffcail/gateway-web/public"
)

type DashboardController struct{}

// RegisterDashboard
func RegisterDashboard(group *gin.RouterGroup) {
	dashboard := &DashboardController{}
	group.GET("/panel_group_data", dashboard.PanelGroupData)
	group.GET("/flow_stat", dashboard.FlowStat)
	group.GET("/service_stat", dashboard.ServiceStat)
}

// PanelGroupData godoc
// @Summary 指标统计
// @Description 指标统计
// @Tags 首页大盘
// @ID /dashboard/panel_group_data
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.PanelGroupDataOutput} "success"
// @Router /dashboard/panel_group_data [get]
func (dashboard *DashboardController) PanelGroupData(c *gin.Context) {
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	serviceInfo := &dao.GatewayServiceInfo{}
	_, serviceNum, err := serviceInfo.PageList(c, tx, &dto.ServiceListInput{PageSize: 1, PageNo: 1})
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	app := &dao.App{}
	_, appNum, err := app.AppList(c, tx, &dto.APPListInput{PageSize: 1, PageNo: 1})
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	counter, err := public.FlowCounterHandler.GetCounter(public.FlowTotal)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	out := &dto.PanelGroupDataOutput{
		ServiceNum:      serviceNum,
		AppNum:          appNum,
		CurrentQPS:      counter.TotalCount,
		TodayRequestNum: counter.Qps,
	}
	middleware.ResponseSuccess(c, out)
}

// FlowStat godoc
// @Summary 流量统计
// @Description 流量统计
// @Tags 首页大盘
// @ID /dashboard/flow_stat
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.DashServiceStatOutput} "success"
// @Router /dashboard/flow_stat [get]
func (dashboard *DashboardController) FlowStat(c *gin.Context) {
	counter, err := public.FlowCounterHandler.GetCounter(public.FlowTotal)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	todayList := []int64{}
	currentTime := time.Now()
	for i := 0; i <= currentTime.Hour(); i++ {
		dateTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), i, 0, 0, 0, lib.TimeLocation)
		hourData, _ := counter.GetHourData(dateTime)
		todayList = append(todayList, hourData)
	}

	yesterdayList := []int64{}
	yesterTime := currentTime.Add(-1 * time.Duration(time.Hour*24))
	for i := 0; i <= 23; i++ {
		dateTime := time.Date(yesterTime.Year(), yesterTime.Month(), yesterTime.Day(), i, 0, 0, 0, lib.TimeLocation)
		hourData, _ := counter.GetHourData(dateTime)
		yesterdayList = append(yesterdayList, hourData)
	}
	middleware.ResponseSuccess(c, &dto.ServiceStatOutput{
		Today:     todayList,
		Yesterday: yesterdayList,
	})
}

// ServiceStat godoc
// @Summary 服务统计
// @Description 服务统计
// @Tags 首页大盘
// @ID /dashboard/service_stat
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.DashServiceStatOutput} "success"
// @Router /dashboard/service_stat [get]
func (dashboard *DashboardController) ServiceStat(c *gin.Context) {
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	serviceInfo := &dao.GatewayServiceInfo{}
	list, err := serviceInfo.GroupByLoadType(c, tx)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	legend := []string{}
	for index, item := range list {
		name, ok := public.LoadTypeMap[item.LoadType]
		if !ok {
			middleware.ResponseError(c, 2003, errors.New("load_type not found"))
			return
		}
		list[index].Name = name
		legend = append(legend, name)
	}
	out := &dto.DashServiceStatOutput{
		Legend: legend,
		Data:   list,
	}
	middleware.ResponseSuccess(c, out)
}
