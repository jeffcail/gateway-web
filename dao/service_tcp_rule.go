package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TcpRule struct {
	ID        int64 `json:"id" gorm:"primary_key"`
	ServiceID int64 `json:"service_id" gorm:"column:service_id"` // 服务id
	Port      int   `json:"port" gorm:"column:port"`             // 端口
}

func (tr *TcpRule) TableName() string {
	return "gateway_service_tcp_rule"
}

// Find
// 查找
func (tr *TcpRule) Find(c *gin.Context, tx *gorm.DB, search *TcpRule) (*TcpRule, error) {
	data := &TcpRule{}
	err := tx.WithContext(c).Where(search).Find(data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Save
// 保存
func (tr *TcpRule) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.WithContext(c).Save(tr).Error; err != nil {
		return err
	}
	return nil
}
