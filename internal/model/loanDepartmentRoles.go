package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
)

type LoanDepartmentRoles struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	DepartmentID int64 `gorm:"column:department_id;type:bigint(20);primary_key" json:"departmentID"`
	RoleID       int64 `gorm:"column:role_id;type:bigint(20);not null" json:"roleID"`
}

// LoanDepartmentRolesColumnNames Whitelist for custom query fields to prevent sql injection attacks
var LoanDepartmentRolesColumnNames = map[string]bool{
	"id":            true,
	"created_at":    true,
	"updated_at":    true,
	"deleted_at":    true,
	"department_id": true,
	"role_id":       true,
}
