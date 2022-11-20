package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GatewayServiceInfo
// 网关基本信息表
type GatewayServiceInfo struct {
	ID          int64  `json:"id" gorm:"column:id;primary_key;AUTO_INCREMENT"`   // 主键id
	LoadType    int    `json:"load_type" gorm:"column:load_type;NOT NULL"`       // 服务类型
	ServiceName string `json:"service_name" gorm:"column:service_name;NOT NULL"` // 服务名称
	ServiceDesc string `json:"service_desc" gorm:"column:service_desc;NOT NULL"` // 服务描述
	CreateAt    string `json:"create_at" gorm:"column:create_at;NOT NULL"`       // 添加时间
	UpdateAt    string `json:"update_at" gorm:"column:update_at;NOT NULL"`       // 更新时间
	IsDelete    int    `json:"is_delete" gorm:"column:is_delete;NOT NULL"`       // 是否删除
}

func (gsi *GatewayServiceInfo) TableName() string {
	return "gateway_service_info"
}

func (gsi *GatewayServiceInfo) Find(c *gin.Context, tx *gorm.DB, search *GatewayServiceInfo) (*GatewayServiceInfo, error) {
	out := &GatewayServiceInfo{}
	err := tx.WithContext(c).Where(search).First(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (gsi *GatewayServiceInfo) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.WithContext(c).Save(gsi).Error
}
