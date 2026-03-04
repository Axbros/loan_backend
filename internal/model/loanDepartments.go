package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
)

type LoanDepartments struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	Name        string `gorm:"column:name;type:varchar(128);not null" json:"name"`
	AdminUserID int64  `gorm:"column:admin_user_id;type:bigint(20)" json:"adminUserID"`
	Status      int    `gorm:"column:status;type:tinyint(4);default:1;not null" json:"status"`
}

// LoanDepartmentsColumnNames Whitelist for custom query fields to prevent sql injection attacks
var LoanDepartmentsColumnNames = map[string]bool{
	"id":            true,
	"created_at":    true,
	"updated_at":    true,
	"deleted_at":    true,
	"name":          true,
	"admin_user_id": true,
	"status":        true,
}
