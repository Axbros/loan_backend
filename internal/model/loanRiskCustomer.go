package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
)

type LoanRiskCustomer struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	LoanBaseinfoID uint64 `gorm:"column:loan_baseinfo_id;type:int(11)" json:"loanBaseinfoID"`
	RiskType       int    `gorm:"column:risk_type;type:tinyint(4)" json:"riskType"`       // 风险类型 -1 黑名单 1 白名单
	RiskReason     string `gorm:"column:risk_reason;type:varchar(255)" json:"riskReason"` // 风险原因
	CreatedBy      uint64 `gorm:"column:created_by;type:int(11)" json:"createdBy"`        // loan_users_id

	LoanBaseinfo *LoanBaseinfo `gorm:"foreignKey:LoanBaseinfoID;references:ID;PRELOAD:false" json:"loanBaseinfo"` // 一对一关联，默认不预加载

	OperateUser *LoanUsers `gorm:"foreignKey:CreatedBy;references:ID;PRELOAD:false" json:"operateUser"`
}

// TableName table name
func (m *LoanRiskCustomer) TableName() string {
	return "loan_risk_customer"
}

// LoanRiskCustomerColumnNames Whitelist for custom query fields to prevent sql injection attacks
var LoanRiskCustomerColumnNames = map[string]bool{
	"id":               true,
	"created_at":       true,
	"updated_at":       true,
	"deleted_at":       true,
	"loan_baseinfo_id": true,
	"risk_type":        true,
	"risk_reason":      true,
	"created_by":       true,
}
