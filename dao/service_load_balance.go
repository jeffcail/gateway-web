package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
)

// GatewayServiceLoadBalance
// 网关负载表
type GatewayServiceLoadBalance struct {
	ID                     int64  `gorm:"column:id;primary_key;AUTO_INCREMENT"`               // 主键id
	ServiceID              int64  `gorm:"column:service_id;default:0;NOT NULL"`               // 服务id
	CheckMethod            string `gorm:"column:check_method;default:0;NOT NULL"`             // 检查方法 0=tcpchk,检测端口是否握手成功
	CheckTimeout           int    `gorm:"column:check_timeout;default:0;NOT NULL"`            // check超时时间,单位s
	CheckInterval          int    `gorm:"column:check_interval;default:0;NOT NULL"`           // 检查间隔, 单位s
	RoundType              int    `gorm:"column:round_type;default:2;NOT NULL"`               // 轮询方式 0=random 1=round-robin 2=weight_round-robin 3=ip_hash
	IPList                 string `gorm:"column:ip_list;NOT NULL"`                            // ip列表
	WeightList             string `gorm:"column:weight_list;NOT NULL"`                        // 权重列表
	ForbidList             string `gorm:"column:forbid_list;NOT NULL"`                        // 禁用ip列表
	UpstreamConnectTimeout int    `gorm:"column:upstream_connect_timeout;default:0;NOT NULL"` // 建立连接超时, 单位s
	UpstreamHeaderTimeout  int    `gorm:"column:upstream_header_timeout;default:0;NOT NULL"`  // 获取header超时, 单位s
	UpstreamIdleTimeout    int    `gorm:"column:upstream_idle_timeout;default:0;NOT NULL"`    // 链接最大空闲时间, 单位s
	UpstreamMaxIdle        int    `gorm:"column:upstream_max_idle;default:0;NOT NULL"`        // 最大空闲链接数
}

func (m *GatewayServiceLoadBalance) TableName() string {
	return "gateway_service_load_balance"
}

// Find
// 查找
func (m *GatewayServiceLoadBalance) Find(c *gin.Context, tx *gorm.DB, search *GatewayServiceLoadBalance) (*GatewayServiceLoadBalance, error) {
	data := &GatewayServiceLoadBalance{}
	err := tx.WithContext(c).Where(search).Find(data).Error
	return data, err
}

// Save
// 保存
func (m *GatewayServiceLoadBalance) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.WithContext(c).Save(m).Error; err != nil {
		return err
	}
	return nil
}

func (m *GatewayServiceLoadBalance) GetIpListByModel() []string {
	return strings.Split(m.IPList, ",")
}
