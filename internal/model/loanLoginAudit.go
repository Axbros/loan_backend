package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
)

type LoanLoginAudit struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	UserID    int64  `gorm:"column:user_id;type:bigint(20)" json:"userID"`
	LoginType string `gorm:"column:login_type;type:varchar(16);not null" json:"loginType"`
	IP        string `gorm:"column:ip;type:varchar(64)" json:"ip"`
	UserAgent string `gorm:"column:user_agent;type:varchar(255)" json:"userAgent"`
	Success   int    `gorm:"column:success;type:tinyint(4);not null" json:"success"`
}

// TableName table name
func (m *LoanLoginAudit) TableName() string {
	return "loan_login_audit"
}

// LoanLoginAuditColumnNames Whitelist for custom query fields to prevent sql injection attacks
var LoanLoginAuditColumnNames = map[string]bool{
	"id":         true,
	"created_at": true,
	"updated_at": true,
	"deleted_at": true,
	"user_id":    true,
	"login_type": true,
	"ip":         true,
	"user_agent": true,
	"success":    true,
}
