package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
)

type LoanRolePermissions struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	RoleID       int64 `gorm:"column:role_id;type:bigint(20);primary_key" json:"roleID"`
	PermissionID int64 `gorm:"column:permission_id;type:bigint(20);not null" json:"permissionID"`
}

// LoanRolePermissionsColumnNames Whitelist for custom query fields to prevent sql injection attacks
var LoanRolePermissionsColumnNames = map[string]bool{
	"id":            true,
	"created_at":    true,
	"updated_at":    true,
	"deleted_at":    true,
	"role_id":       true,
	"permission_id": true,
}
