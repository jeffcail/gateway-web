package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GatewayServiceAccessControl
// 网关权限控制表
type GatewayServiceAccessControl struct {
	ID                int64  `gorm:"column:id;primary_key;AUTO_INCREMENT"` // 自增主键
	ServiceID         int64  `gorm:"column:service_id"`                    // 服务id
	OpenAuth          int    `gorm:"column:open_auth"`                     // 是否开启权限 1=开启
	BlackList         string `gorm:"column:black_list"`                    // 黑名单ip
	WhiteList         string `gorm:"column:white_list"`                    // 白名单ip
	WhiteHostName     string `gorm:"column:white_host_name"`               // 白名单主机
	ClientipFlowLimit int    `gorm:"column:clientip_flow_limit"`           // 客户端ip限流
	ServiceFlowLimit  int    `gorm:"column:service_flow_limit"`            // 服务端限流
}

func (m *GatewayServiceAccessControl) TableName() string {
	return "gateway_service_access_control"
}

// Find
// 查找
func (m *GatewayServiceAccessControl) Find(c *gin.Context, tx *gorm.DB, search *GatewayServiceAccessControl) (
	*GatewayServiceAccessControl, error) {
	data := &GatewayServiceAccessControl{}
	if err := tx.WithContext(c).Where(search).Find(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Save
// 保存
func (m *GatewayServiceAccessControl) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.WithContext(c).Save(m).Error; err != nil {
		return err
	}
	return nil
}
