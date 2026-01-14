package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
)

type LoanRoleDepartments struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	RoleID       int64 `gorm:"column:role_id;type:bigint(20);primary_key" json:"roleID"`
	DepartmentID int64 `gorm:"column:department_id;type:bigint(20);not null" json:"departmentID"`
}

// LoanRoleDepartmentsColumnNames Whitelist for custom query fields to prevent sql injection attacks
var LoanRoleDepartmentsColumnNames = map[string]bool{
	"id":            true,
	"created_at":    true,
	"updated_at":    true,
	"deleted_at":    true,
	"role_id":       true,
	"department_id": true,
}
