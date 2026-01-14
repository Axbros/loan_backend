package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
)

type LoanPermissions struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	Code     string `gorm:"column:code;type:varchar(128);not null" json:"code"`
	Name     string `gorm:"column:name;type:varchar(128);not null" json:"name"`
	Type     string `gorm:"column:type;type:varchar(16)" json:"type"`
	Resource string `gorm:"column:resource;type:varchar(255)" json:"resource"`
}

// LoanPermissionsColumnNames Whitelist for custom query fields to prevent sql injection attacks
var LoanPermissionsColumnNames = map[string]bool{
	"id":         true,
	"created_at": true,
	"updated_at": true,
	"deleted_at": true,
	"code":       true,
	"name":       true,
	"type":       true,
	"resource":   true,
}
