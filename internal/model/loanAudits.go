package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
)

// LoanAudits 申请审核记录表(审核时间即 created_at)
type LoanAudits struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	BaseinfoID    uint64 `gorm:"column:baseinfo_id;type:int(11);not null" json:"baseinfoID"`           // 关联申请单 loan_baseinfo.id
	AuditResult   int    `gorm:"column:audit_result;type:tinyint(4);not null" json:"auditResult"`      // 审核结果：1通过 -1拒绝
	AuditComment  string `gorm:"column:audit_comment;type:varchar(255)" json:"auditComment"`           // 审核备注/原因
	AuditorUserID uint64 `gorm:"column:auditor_user_id;type:bigint(20);not null" json:"auditorUserID"` // 审核人员(loan_users.id)
	AuditType     int    `gorm:"column:audit_type;type:tinyint(4);not null" json:"auditType"`          // '审核类型(初审0、放款审核1、回款审核2)',
}

// LoanAuditsColumnNames Whitelist for custom query fields to prevent sql injection attacks
var LoanAuditsColumnNames = map[string]bool{
	"id":              true,
	"created_at":      true,
	"updated_at":      true,
	"deleted_at":      true,
	"baseinfo_id":     true,
	"audit_type":      true,
	"audit_result":    true,
	"audit_comment":   true,
	"auditor_user_id": true,
}
