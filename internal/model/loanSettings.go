package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
)

type LoanSettings struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	Name   string `gorm:"column:name;type:varchar(255)" json:"name"`
	Value  string `gorm:"column:value;type:varchar(255)" json:"value"`
	Remark string `gorm:"column:remark;type:varchar(255)" json:"remark"`
}

// LoanSettingsColumnNames Whitelist for custom query fields to prevent sql injection attacks
var LoanSettingsColumnNames = map[string]bool{
	"id":         true,
	"created_at": true,
	"updated_at": true,
	"deleted_at": true,
	"name":       true,
	"value":      true,
	"remark":     true,
}
