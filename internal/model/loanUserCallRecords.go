package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
	"time"
)

// LoanUserCallRecords 通话记录采集表(匿名表单采集，与loan_baseinfo关联)
type LoanUserCallRecords struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	BaseinfoID      int        `gorm:"column:baseinfo_id;type:int(11);not null" json:"baseinfoID"`                     // 关联 loan_baseinfo.id
	CallType        int        `gorm:"column:call_type;type:tinyint(4);not null" json:"callType"`                      // 通话类型：1呼入 2呼出 3未接 4拒接(按采集端定义)
	PhoneNumber     string     `gorm:"column:phone_number;type:varchar(32)" json:"phoneNumber"`                        // 对端号码/电话
	PhoneNormalized string     `gorm:"column:phone_normalized;type:varchar(32)" json:"phoneNormalized"`                // 标准化号码(去空格/国家码等，可选)
	CallTime        *time.Time `gorm:"column:call_time;type:datetime" json:"callTime"`                                 // 通话开始时间(手机侧时间)
	DurationSeconds int        `gorm:"column:duration_seconds;type:int(11);default:0;not null" json:"durationSeconds"` // 通话时长(秒，未接/拒接一般为0)
	CallHash        string     `gorm:"column:call_hash;type:char(64)" json:"callHash"`                                 // 记录去重哈希(如 sha256(type+phone+call_time+duration)，可选)
}

// LoanUserCallRecordsColumnNames Whitelist for custom query fields to prevent sql injection attacks
var LoanUserCallRecordsColumnNames = map[string]bool{
	"id":               true,
	"created_at":       true,
	"updated_at":       true,
	"deleted_at":       true,
	"baseinfo_id":      true,
	"call_type":        true,
	"phone_number":     true,
	"phone_normalized": true,
	"call_time":        true,
	"duration_seconds": true,
	"call_hash":        true,
}
