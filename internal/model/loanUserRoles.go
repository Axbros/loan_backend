package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
)

type LoanUserRoles struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	UserID int64 `gorm:"column:user_id;type:bigint(20);primary_key" json:"userID"`
	RoleID int64 `gorm:"column:role_id;type:bigint(20);not null" json:"roleID"`
}

// LoanUserRolesColumnNames Whitelist for custom query fields to prevent sql injection attacks
var LoanUserRolesColumnNames = map[string]bool{
	"id":         true,
	"created_at": true,
	"updated_at": true,
	"deleted_at": true,
	"user_id":    true,
	"role_id":    true,
}
