package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
)

type LoanRoles struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	Code      string `gorm:"column:code;type:varchar(64);not null" json:"code"`
	Name      string `gorm:"column:name;type:varchar(128);not null" json:"name"`
	DataScope string `gorm:"column:data_scope;type:varchar(32);default:DEPT;not null" json:"dataScope"`
	Status    int    `gorm:"column:status;type:tinyint(4);default:1;not null" json:"status"`
}

// LoanRolesColumnNames Whitelist for custom query fields to prevent sql injection attacks
var LoanRolesColumnNames = map[string]bool{
	"id":         true,
	"created_at": true,
	"updated_at": true,
	"deleted_at": true,
	"code":       true,
	"name":       true,
	"data_scope": true,
	"status":     true,
}
