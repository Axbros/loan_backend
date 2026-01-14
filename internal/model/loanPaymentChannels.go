package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
	"github.com/shopspring/decimal"
)

// LoanPaymentChannels 统一支付渠道配置表(支持代付放款+代收回款，可禁用/启用，含手续费/限额/结算周期)
type LoanPaymentChannels struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	Code             string           `gorm:"column:code;type:varchar(64);not null" json:"code"`                       // 渠道编码(唯一，如 BANK_A、WALLET_X)
	Name             string           `gorm:"column:name;type:varchar(128);not null" json:"name"`                      // 渠道名称
	MerchantNo       string           `gorm:"column:merchant_no;type:varchar(128);not null" json:"merchantNo"`         // 商户号/商户ID(该渠道分配给平台的商户标识)
	Status           int              `gorm:"column:status;type:tinyint(4);default:1;not null" json:"status"`          // 渠道状态：1启用 0禁用
	CanPayout        int              `gorm:"column:can_payout;type:tinyint(4);default:1;not null" json:"canPayout"`   // 是否支持代付/放款：1是 0否
	CanCollect       int              `gorm:"column:can_collect;type:tinyint(4);default:1;not null" json:"canCollect"` // 是否支持代收/回款：1是 0否
	PayoutFeeRate    *decimal.Decimal `gorm:"column:payout_fee_rate;type:decimal(10,6)" json:"payoutFeeRate"`          // 代付手续费率(如0.003500=0.35%)
	PayoutFeeFixed   int              `gorm:"column:payout_fee_fixed;type:int(11)" json:"payoutFeeFixed"`              // 代付固定手续费(分，若不用可为空)
	CollectFeeRate   *decimal.Decimal `gorm:"column:collect_fee_rate;type:decimal(10,6)" json:"collectFeeRate"`        // 代收手续费率
	CollectFeeFixed  int              `gorm:"column:collect_fee_fixed;type:int(11)" json:"collectFeeFixed"`            // 代收固定手续费(分)
	CollectMinAmount int              `gorm:"column:collect_min_amount;type:int(11)" json:"collectMinAmount"`          // 最小代收金额(分)
	CollectMaxAmount int              `gorm:"column:collect_max_amount;type:int(11)" json:"collectMaxAmount"`          // 最大代收金额(分)
	PayoutMinAmount  int              `gorm:"column:payout_min_amount;type:int(11)" json:"payoutMinAmount"`            // 最小代付金额(分)
	PayoutMaxAmount  int              `gorm:"column:payout_max_amount;type:int(11)" json:"payoutMaxAmount"`            // 最大代付金额(分)
	SettlementCycle  string           `gorm:"column:settlement_cycle;type:varchar(32)" json:"settlementCycle"`         // 结算周期(如 T0/T1/D1/W1/M1，可按你们渠道定义)
	SettlementDesc   string           `gorm:"column:settlement_desc;type:varchar(255)" json:"settlementDesc"`          // 结算说明/备注
}

// LoanPaymentChannelsColumnNames Whitelist for custom query fields to prevent sql injection attacks
var LoanPaymentChannelsColumnNames = map[string]bool{
	"id":                 true,
	"created_at":         true,
	"updated_at":         true,
	"deleted_at":         true,
	"code":               true,
	"name":               true,
	"merchant_no":        true,
	"status":             true,
	"can_payout":         true,
	"can_collect":        true,
	"payout_fee_rate":    true,
	"payout_fee_fixed":   true,
	"collect_fee_rate":   true,
	"collect_fee_fixed":  true,
	"collect_min_amount": true,
	"collect_max_amount": true,
	"payout_min_amount":  true,
	"payout_max_amount":  true,
	"settlement_cycle":   true,
	"settlement_desc":    true,
}
