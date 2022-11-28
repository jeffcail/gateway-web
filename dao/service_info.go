package dao

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jeffcail/gateway-web/dto"
	"gorm.io/gorm"
)

// GatewayServiceInfo
// 网关基本信息表
type GatewayServiceInfo struct {
	ID          int64     `json:"id" gorm:"column:id;primary_key;AUTO_INCREMENT"`   // 主键id
	LoadType    int       `json:"load_type" gorm:"column:load_type;NOT NULL"`       // 服务类型
	ServiceName string    `json:"service_name" gorm:"column:service_name;NOT NULL"` // 服务名称
	ServiceDesc string    `json:"service_desc" gorm:"column:service_desc;NOT NULL"` // 服务描述
	CreatedAt   time.Time `json:"create_at" gorm:"column:create_at;NOT NULL"`       // 添加时间
	UpdatedAt   time.Time `json:"update_at" gorm:"column:update_at;NOT NULL"`       // 更新时间
	IsDelete    int       `json:"is_delete" gorm:"column:is_delete;NOT NULL"`       // 是否删除
}

func (gsi *GatewayServiceInfo) TableName() string {
	return "gateway_service_info"
}

func (gsi *GatewayServiceInfo) Find(c *gin.Context, tx *gorm.DB, search *GatewayServiceInfo) (*GatewayServiceInfo, error) {
	out := &GatewayServiceInfo{}
	err := tx.WithContext(c).Where(search).First(out).Error
	return out, err
}

func (gsi *GatewayServiceInfo) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.WithContext(c).Save(gsi).Error
}

func (t *GatewayServiceInfo) ServiceDetail(c *gin.Context, tx *gorm.DB, search *GatewayServiceInfo) (
	*ServiceDetail, error) {
	if search.ID == 0 {
		info, err := t.Find(c, tx, search)
		if err != nil {
			return nil, err
		}
		search = info
	}
	httpRule := &GatewayServiceHttpRule{ServiceID: search.ID}
	httpRule, err := httpRule.Find(c, tx, httpRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	tcpRule := &TcpRule{ServiceID: search.ID}
	tcpRule, err = tcpRule.Find(c, tx, tcpRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	grpcRule := &GrpcRule{ServiceID: search.ID}
	grpcRule, err = grpcRule.Find(c, tx, grpcRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	accessControl := &GatewayServiceAccessControl{ServiceID: search.ID}
	accessControl, err = accessControl.Find(c, tx, accessControl)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	loadBalance := &GatewayServiceLoadBalance{ServiceID: search.ID}
	accessControl, err = accessControl.Find(c, tx, accessControl)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	detail := &ServiceDetail{
		Info:          search,
		HTTPRule:      httpRule,
		TcpRule:       tcpRule,
		GrpcRule:      grpcRule,
		LoadBalance:   loadBalance,
		AccessControl: accessControl,
	}
	return detail, nil
}

func (gsi *GatewayServiceInfo) PageList(c *gin.Context, tx *gorm.DB, param *dto.ServiceListInput) (
	[]GatewayServiceInfo, int64, error) {
	total := int64(0)
	list := []GatewayServiceInfo{}
	offset := (param.PageNo - 1) * param.PageSize

	query := tx.WithContext(c)
	query = query.Table(gsi.TableName()).Where("is_delete=0")
	if param.Info != "" {
		query = query.Where("(service_name like ? or service_desc like ?)", "%"+param.Info+"%", "%"+param.Info+"%")
	}
	if err := query.Limit(param.PageSize).Offset(offset).Order("id desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	query.Limit(param.PageSize).Offset(offset).Count(&total)
	return list, total, nil
}
