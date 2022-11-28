package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GatewayServiceHttpRule
// http网关路由匹配表
type GatewayServiceHttpRule struct {
	ID             int64  `gorm:"column:id;primary_key;AUTO_INCREMENT"`     // 主键id
	ServiceID      int64  `gorm:"column:service_id;NOT NULL"`               // 服务id
	RuleType       int    `gorm:"column:rule_type;default:0;NOT NULL"`      // 匹配类型 0=url前缀url_prefix 1=域名domain
	Rule           string `gorm:"column:rule;NOT NULL"`                     // type=domain表示域名，type=url_prefix时表示url前缀
	NeedHttps      int    `gorm:"column:need_https;default:0;NOT NULL"`     // 支持https 1=支持
	NeedStripUri   int    `gorm:"column:need_strip_uri;default:0;NOT NULL"` // 启用strip_uri 1=启用
	NeedWebsocket  int    `gorm:"column:need_websocket;default:0;NOT NULL"` // 是否支持websocket 1=支持
	UrlRewrite     string `gorm:"column:url_rewrite;NOT NULL"`              // url重写功能 格式：^/gatekeeper/test_service(.*) $1 多个逗号间隔
	HeaderTransfor string `gorm:"column:header_transfor"`                   // header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue 多个逗号间隔
}

func (m *GatewayServiceHttpRule) TableName() string {
	return "gateway_service_http_rule"
}

// Find
// 查找
func (m *GatewayServiceHttpRule) Find(c *gin.Context, tx *gorm.DB, search *GatewayServiceHttpRule) (
	*GatewayServiceHttpRule, error) {
	data := &GatewayServiceHttpRule{}
	err := tx.WithContext(c).Where(search).Find(data).Error
	return data, err
}

// Save
// 保存
func (m *GatewayServiceHttpRule) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.WithContext(c).Save(m).Error; err != nil {
		return err
	}
	return nil
}
