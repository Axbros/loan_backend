package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
	"time"
)

// LoanRepaymentSchedules 还款计划表(支持单期/分期，逾期/已还通过状态体现)
type LoanRepaymentSchedules struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	DisbursementID int64      `gorm:"column:disbursement_id;type:bigint(20);not null" json:"disbursementID"`       // 关联放款单 loan_disbursements.id
	InstallmentNo  int        `gorm:"column:installment_no;type:int(11);not null" json:"installmentNo"`            // 期次(从1开始)
	DueDate        *time.Time `gorm:"column:due_date;type:date;not null" json:"dueDate"`                           // 应还日期
	PrincipalDue   int64      `gorm:"column:principal_due;type:bigint(20);default:0;not null" json:"principalDue"` // 应还本金(建议统一单位：分)
	InterestDue    int64      `gorm:"column:interest_due;type:bigint(20);default:0;not null" json:"interestDue"`   // 应还利息(分)
	FeeDue         int64      `gorm:"column:fee_due;type:bigint(20);default:0;not null" json:"feeDue"`             // 应还费用(分)
	PenaltyDue     int        `gorm:"column:penalty_due;type:int(11);default:0;not null" json:"penaltyDue"`        // 应还罚息(分，逾期产生)
	TotalDue       int64      `gorm:"column:total_due;type:bigint(20);not null" json:"totalDue"`                   // 本期应还总额=本金+利息+费用+罚息(分)
	PaidPrincipal  int        `gorm:"column:paid_principal;type:int(11);default:0;not null" json:"paidPrincipal"`  // 已还本金(分)
	PaidInterest   int        `gorm:"column:paid_interest;type:int(11);default:0;not null" json:"paidInterest"`    // 已还利息(分)
	PaidFee        int        `gorm:"column:paid_fee;type:int(11);default:0;not null" json:"paidFee"`              // 已还费用(分)
	PaidPenalty    int        `gorm:"column:paid_penalty;type:int(11);default:0;not null" json:"paidPenalty"`      // 已还罚息(分)
	PaidTotal      int        `gorm:"column:paid_total;type:int(11);default:0;not null" json:"paidTotal"`          // 已还总额(分)
	Status         int        `gorm:"column:status;type:tinyint(4);default:0;not null" json:"status"`              // 期次状态：0未还清 1已还清 2逾期
	LastPaidAt     *time.Time `gorm:"column:last_paid_at;type:datetime" json:"lastPaidAt"`                         // 最近一次还款时间
	SettledAt      *time.Time `gorm:"column:settled_at;type:datetime" json:"settledAt"`                            // 结清时间(本期还清时)
}

// LoanRepaymentSchedulesColumnNames Whitelist for custom query fields to prevent sql injection attacks
var LoanRepaymentSchedulesColumnNames = map[string]bool{
	"id":              true,
	"created_at":      true,
	"updated_at":      true,
	"deleted_at":      true,
	"disbursement_id": true,
	"installment_no":  true,
	"due_date":        true,
	"principal_due":   true,
	"interest_due":    true,
	"fee_due":         true,
	"penalty_due":     true,
	"total_due":       true,
	"paid_principal":  true,
	"paid_interest":   true,
	"paid_fee":        true,
	"paid_penalty":    true,
	"paid_total":      true,
	"status":          true,
	"last_paid_at":    true,
	"settled_at":      true,
}
