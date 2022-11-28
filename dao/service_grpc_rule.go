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

func (t *GrpcRule) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.WithContext(c).Save(t).Error; err != nil {
		return err
	}
	return nil
}

func (t *GrpcRule) ListByServiceID(c *gin.Context, tx *gorm.DB, serviceID int64) ([]GrpcRule, int64, error) {
	var list []GrpcRule
	var count int64
	query := tx.WithContext(c)
	query = query.Table(t.TableName()).Select("*")
	query = query.Where("service_id = ?", serviceID)
	err := query.Order("id desc").Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	errCount := query.Count(&count).Error
	if errCount != nil {
		return nil, 0, err
	}
	return list, count, nil
}
