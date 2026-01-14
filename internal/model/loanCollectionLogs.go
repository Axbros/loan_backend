package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
	"time"
)

// LoanCollectionLogs 催收跟进记录表(一条任务可多次记录沟通内容/承诺/计划)
type LoanCollectionLogs struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	CaseID          int64      `gorm:"column:case_id;type:bigint(20);not null" json:"caseID"`                    // 关联催收任务 loan_collection_cases.id
	CollectorUserID int64      `gorm:"column:collector_user_id;type:bigint(20);not null" json:"collectorUserID"` // 催收人员 loan_users.id
	ActionType      string     `gorm:"column:action_type;type:varchar(32)" json:"actionType"`                    // 动作类型(如 CALL/SMS/VISIT/OTHER，可选)
	Content         string     `gorm:"column:content;type:varchar(500);not null" json:"content"`                 // 跟进内容/备注(例如用户承诺3天内还款)
	NextFollowUpAt  *time.Time `gorm:"column:next_follow_up_at;type:datetime" json:"nextFollowUpAt"`             // 下次跟进时间(可选)
}

// LoanCollectionLogsColumnNames Whitelist for custom query fields to prevent sql injection attacks
var LoanCollectionLogsColumnNames = map[string]bool{
	"id":                true,
	"created_at":        true,
	"updated_at":        true,
	"deleted_at":        true,
	"case_id":           true,
	"collector_user_id": true,
	"action_type":       true,
	"content":           true,
	"next_follow_up_at": true,
}
