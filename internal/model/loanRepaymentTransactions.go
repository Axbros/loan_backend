package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
	"time"
)

// LoanRepaymentTransactions 回款流水表(记录每次实际回款，支持分期/部分还款，含回款渠道与订单号)
type LoanRepaymentTransactions struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	ScheduleID       int64      `gorm:"column:schedule_id;type:bigint(20)" json:"scheduleID"`                         // 关联期次 loan_repayment_schedules.id(可空：先入账后分配/未分期)
	CollectChannelID int64      `gorm:"column:collect_channel_id;type:bigint(20)" json:"collectChannelID"`            // 回款渠道(代收) loan_payment_channels.id
	CollectOrderNo   string     `gorm:"column:collect_order_no;type:varchar(128)" json:"collectOrderNo"`              // 回款订单号/三方代收单号(商户单号)
	PayRef           string     `gorm:"column:pay_ref;type:varchar(128)" json:"payRef"`                               // 支付渠道流水号/交易号(三方transaction id)
	PayAmount        int        `gorm:"column:pay_amount;type:int(11);not null" json:"payAmount"`                     // 本次回款金额(分)
	PayMethod        string     `gorm:"column:pay_method;type:varchar(32)" json:"payMethod"`                          // 回款方式(如 BANK_TRANSFER/CARD/WALLET/CASH)
	PaidAt           *time.Time `gorm:"column:paid_at;type:datetime;not null" json:"paidAt"`                          // 回款时间(交易成功时间)
	AllocPrincipal   int        `gorm:"column:alloc_principal;type:int(11);default:0;not null" json:"allocPrincipal"` // 本次分配到本金(分)
	AllocInterest    int        `gorm:"column:alloc_interest;type:int(11);default:0;not null" json:"allocInterest"`   // 本次分配到利息(分)
	AllocFee         int        `gorm:"column:alloc_fee;type:int(11);default:0;not null" json:"allocFee"`             // 本次分配到费用(分)
	AllocPenalty     int        `gorm:"column:alloc_penalty;type:int(11);default:0;not null" json:"allocPenalty"`     // 本次分配到罚息(分)
	Status           int        `gorm:"column:status;type:tinyint(4);default:1;not null" json:"status"`               // 流水状态：1成功 0失败 2冲正/撤销
	VoucherFileName  string     `gorm:"column:voucher_file_name;type:varchar(64)" json:"voucherFileName"`
	Remark           string     `gorm:"column:remark;type:varchar(255)" json:"remark"` // 备注
	CreatedBy        uint64     `gorm:"column:created_by;type:bigint(20)" json:"createdBy"`
}

// LoanRepaymentTransactionsColumnNames Whitelist for custom query fields to prevent sql injection attacks
var LoanRepaymentTransactionsColumnNames = map[string]bool{
	"id":                 true,
	"created_at":         true,
	"updated_at":         true,
	"deleted_at":         true,
	"disbursement_id":    true,
	"schedule_id":        true,
	"collect_channel_id": true,
	"collect_order_no":   true,
	"pay_ref":            true,
	"pay_amount":         true,
	"pay_method":         true,
	"paid_at":            true,
	"alloc_principal":    true,
	"alloc_interest":     true,
	"alloc_fee":          true,
	"alloc_penalty":      true,
	"status":             true,
	"remark":             true,
}
