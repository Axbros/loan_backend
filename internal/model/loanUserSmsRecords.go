package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
	"time"
)

// LoanUserSmsRecords 短信记录采集表(匿名表单采集，与loan_baseinfo关联)
type LoanUserSmsRecords struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	BaseinfoID int        `gorm:"column:baseinfo_id;type:int(11);not null" json:"baseinfoID"` // 关联 loan_baseinfo.id
	Direction  int        `gorm:"column:direction;type:tinyint(4);not null" json:"direction"` // 短信方向：1收(inbox) 2发(sent)
	Address    string     `gorm:"column:address;type:varchar(64)" json:"address"`             // 对端号码/短码/发件人(如银行短码)
	SmsTime    *time.Time `gorm:"column:sms_time;type:datetime" json:"smsTime"`               // 短信时间(手机侧时间)
	Body       string     `gorm:"column:body;type:text" json:"body"`                          // 短信内容(可选，敏感数据请注意合规)
	BodyHash   string     `gorm:"column:body_hash;type:char(64)" json:"bodyHash"`             // 短信内容哈希(用于去重/审计，可选)
}

// LoanUserSmsRecordsColumnNames Whitelist for custom query fields to prevent sql injection attacks
var LoanUserSmsRecordsColumnNames = map[string]bool{
	"id":          true,
	"created_at":  true,
	"updated_at":  true,
	"deleted_at":  true,
	"baseinfo_id": true,
	"direction":   true,
	"address":     true,
	"sms_time":    true,
	"body":        true,
	"body_hash":   true,
}
