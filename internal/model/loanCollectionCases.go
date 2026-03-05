package model

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm"
)

// LoanCollectionCases 催收任务表(管理员批量分配逾期任务给催收人员，催收人员完成并备注)
type LoanCollectionCases struct {
	sgorm.Model      `gorm:"embedded"` // embed id and time
	ScheduleID       uint64            `gorm:"column:schedule_id;type:bigint(20)" json:"scheduleID"`                        // 关联逾期期次 loan_repayment_schedules.id(按期催收可用，整单催收可为空)
	CollectorUserID  uint64            `gorm:"column:collector_user_id;type:bigint(20);not null" json:"collectorUserID"`    // 催收人员 loan_users.id
	AssignedByUserID uint64            `gorm:"column:assigned_by_user_id;type:bigint(20);not null" json:"assignedByUserID"` // 分配人(管理员) loan_users.id
	Priority         int               `gorm:"column:priority;type:tinyint(4);default:2;not null" json:"priority"`          // 优先级：1高 2中 3低
	Status           int               `gorm:"column:status;type:tinyint(4);default:0;not null" json:"status"`              // 任务状态：0待处理 1跟进中 2已完成 3已取消
	CompletedAt      *time.Time        `gorm:"column:completed_at;type:datetime" json:"completedAt"`                        // 完成时间(点击完成时)
	CompletedNote    string            `gorm:"column:completed_note;type:varchar(255)" json:"completedNote"`                // 完成备注(例如用户承诺X天内还款)
}

// LoanCollectionCasesColumnNames Whitelist for custom query fields to prevent sql injection attacks
var LoanCollectionCasesColumnNames = map[string]bool{
	"id":                  true,
	"created_at":          true,
	"updated_at":          true,
	"deleted_at":          true,
	"schedule_id":         true,
	"collector_user_id":   true,
	"assigned_by_user_id": true,
	"priority":            true,
	"status":              true,
	"completed_at":        true,
	"completed_note":      true,
}
