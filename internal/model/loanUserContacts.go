package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
)

// LoanUserContacts 用户通讯录采集表(匿名表单采集，与loan_baseinfo关联)
type LoanUserContacts struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	BaseinfoID  int    `gorm:"column:baseinfo_id;type:int(11);not null" json:"baseinfoID"` // 关联 loan_baseinfo.id
	ContactName string `gorm:"column:contact_name;type:varchar(128)" json:"contactName"`   // 联系人姓名
	PhoneNumber string `gorm:"column:phone_number;type:varchar(32)" json:"phoneNumber"`    // 联系人手机号/电话
	ContactHash string `gorm:"column:contact_hash;type:char(64)" json:"contactHash"`       // 联系人去重哈希(如 sha256(name+phone_normalized))
}

// LoanUserContactsColumnNames Whitelist for custom query fields to prevent sql injection attacks
var LoanUserContactsColumnNames = map[string]bool{
	"id":           true,
	"created_at":   true,
	"updated_at":   true,
	"deleted_at":   true,
	"baseinfo_id":  true,
	"contact_name": true,
	"phone_number": true,
	"contact_hash": true,
}
