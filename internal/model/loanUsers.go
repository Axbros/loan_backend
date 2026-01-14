package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
)

type LoanUsers struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	Username     string `gorm:"column:username;type:varchar(64);not null" json:"username"`
	PasswordHash string `gorm:"column:password_hash;type:varchar(255);not null" json:"passwordHash"`
	DepartmentID int64  `gorm:"column:department_id;type:bigint(20);not null" json:"departmentID"`
	MfaEnabled   int    `gorm:"column:mfa_enabled;type:tinyint(4);default:0;not null" json:"mfaEnabled"`
	MfaRequired  int    `gorm:"column:mfa_required;type:tinyint(4);default:0;not null" json:"mfaRequired"`
	Status       int    `gorm:"column:status;type:tinyint(4);default:1;not null" json:"status"`
	ShareCode    string `gorm:"column:share_code;type:varchar(32)" json:"shareCode"` // 分享邀请码(用于生成分享链接，建议唯一)
}

// LoanUsersColumnNames Whitelist for custom query fields to prevent sql injection attacks
var LoanUsersColumnNames = map[string]bool{
	"id":            true,
	"created_at":    true,
	"updated_at":    true,
	"deleted_at":    true,
	"username":      true,
	"password_hash": true,
	"department_id": true,
	"mfa_enabled":   true,
	"mfa_required":  true,
	"status":        true,
	"share_code":    true,
}
