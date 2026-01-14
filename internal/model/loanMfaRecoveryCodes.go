package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
	"time"
)

type LoanMfaRecoveryCodes struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	UserID   int64      `gorm:"column:user_id;type:bigint(20);not null" json:"userID"`
	CodeHash string     `gorm:"column:code_hash;type:varbinary(64);not null" json:"codeHash"`
	UsedAt   *time.Time `gorm:"column:used_at;type:datetime" json:"usedAt"`
}

// LoanMfaRecoveryCodesColumnNames Whitelist for custom query fields to prevent sql injection attacks
var LoanMfaRecoveryCodesColumnNames = map[string]bool{
	"id":         true,
	"created_at": true,
	"updated_at": true,
	"deleted_at": true,
	"user_id":    true,
	"code_hash":  true,
	"used_at":    true,
}
