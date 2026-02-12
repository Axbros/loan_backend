package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
	"github.com/shopspring/decimal"
	"time"
)

// LoanDisbursements 放款单/待放款任务表(审核通过后生成，状态待放款->已放款)
type LoanDisbursements struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	BaseinfoID           uint64           `gorm:"column:baseinfo_id;type:int(11);not null" json:"baseinfoID"`                 // 关联申请单 loan_baseinfo.id
	DisburseAmount       *decimal.Decimal `gorm:"column:disburse_amount;type:decimal(10,6);not null" json:"disburseAmount"`   // 放款金额(单位按你的系统：元/分，建议统一)
	NetAmount            *decimal.Decimal `gorm:"column:net_amount;type:decimal(10,6);not null" json:"netAmount"`             // 到账金额(扣除费用后实际到账)
	Status               int              `gorm:"column:status;type:tinyint(4);default:0;not null" json:"status"`             // 放款状态：0待放款 1已放款
	SourceReferrerUserID int64            `gorm:"column:source_referrer_user_id;type:bigint(20)" json:"sourceReferrerUserID"` // 用户来源(分享人 loan_users.id，冗余快照，便于查询)
	AuditorUserID        uint64           `gorm:"column:auditor_user_id;type:bigint(20)" json:"auditorUserID"`                // 审核人员(loan_users.id)
	AuditedAt            *time.Time       `gorm:"column:audited_at;type:datetime" json:"auditedAt"`                           // 审核通过时间
	PayoutChannelID      uint64           `gorm:"column:payout_channel_id;type:bigint(20)" json:"payoutChannelID"`            // 放款渠道(代付) loan_payment_channels.id
	PayoutOrderNo        string           `gorm:"column:payout_order_no;type:varchar(128)" json:"payoutOrderNo"`              // 放款订单号/三方代付单号
	DisbursedAt          *time.Time       `gorm:"column:disbursed_at;type:datetime" json:"disbursedAt"`                       // 放款时间
}

// LoanDisbursementsColumnNames Whitelist for custom query fields to prevent sql injection attacks
var LoanDisbursementsColumnNames = map[string]bool{
	"id":                      true,
	"created_at":              true,
	"updated_at":              true,
	"deleted_at":              true,
	"baseinfo_id":             true,
	"disburse_amount":         true,
	"net_amount":              true,
	"status":                  true,
	"source_referrer_user_id": true,
	"auditor_user_id":         true,
	"audited_at":              true,
	"payout_channel_id":       true,
	"payout_order_no":         true,
	"disbursed_at":            true,
}
