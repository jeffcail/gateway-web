package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GrpcRule struct {
	ID             int64  `json:"id" gorm:"primary_key"`
	ServiceID      int64  `json:"service_id" gorm:"column:service_id"`
	Port           int    `json:"port" gorm:"column:port"`
	HeaderTransfor string `json:"header_transfor" gorm:"column:header_transfor"`
}

func (gr *GrpcRule) TableName() string {
	return "gateway_service_grpc_rule"
}

func (gr *GrpcRule) Find(c *gin.Context, tx *gorm.DB, search *GrpcRule) (*GrpcRule, error) {
	data := &GrpcRule{}
	err := tx.WithContext(c).Where(search).Find(data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}
