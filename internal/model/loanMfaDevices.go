package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
	"time"
)

type LoanMfaDevices struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	UserID     int64      `gorm:"column:user_id;type:bigint(20);not null" json:"userID"`
	Type       string     `gorm:"column:type;type:varchar(16);not null" json:"type"`
	Name       string     `gorm:"column:name;type:varchar(64);not null" json:"name"`
	SecretEnc  string     `gorm:"column:secret_enc;type:varbinary(255)" json:"secretEnc"`
	IsPrimary  int        `gorm:"column:is_primary;type:tinyint(4);default:1;not null" json:"isPrimary"`
	Status     int        `gorm:"column:status;type:tinyint(4);default:1;not null" json:"status"`
	LastUsedAt *time.Time `gorm:"column:last_used_at;type:datetime" json:"lastUsedAt"`
}

// LoanMfaDevicesColumnNames Whitelist for custom query fields to prevent sql injection attacks
var LoanMfaDevicesColumnNames = map[string]bool{
	"id":           true,
	"created_at":   true,
	"updated_at":   true,
	"deleted_at":   true,
	"user_id":      true,
	"type":         true,
	"name":         true,
	"secret_enc":   true,
	"is_primary":   true,
	"status":       true,
	"last_used_at": true,
}
