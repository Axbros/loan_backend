package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateLoanRepaymentTransactionsRequest request params
type CreateLoanRepaymentTransactionsRequest struct {
	DisbursementID   int64      `json:"disbursementID" binding:""`   // 关联放款单 loan_disbursements.id
	ScheduleID       int64      `json:"scheduleID" binding:""`       // 关联期次 loan_repayment_schedules.id(可空：先入账后分配/未分期)
	CollectChannelID int64      `json:"collectChannelID" binding:""` // 回款渠道(代收) loan_payment_channels.id
	CollectOrderNo   string     `json:"collectOrderNo" binding:""`   // 回款订单号/三方代收单号(商户单号)
	PayRef           string     `json:"payRef" binding:""`           // 支付渠道流水号/交易号(三方transaction id)
	PayAmount        int        `json:"payAmount" binding:""`        // 本次回款金额(分)
	PayMethod        string     `json:"payMethod" binding:""`        // 回款方式(如 BANK_TRANSFER/CARD/WALLET/CASH)
	PaidAt           *time.Time `json:"paidAt" binding:""`           // 回款时间(交易成功时间)
	AllocPrincipal   int        `json:"allocPrincipal" binding:""`   // 本次分配到本金(分)
	AllocInterest    int        `json:"allocInterest" binding:""`    // 本次分配到利息(分)
	AllocFee         int        `json:"allocFee" binding:""`         // 本次分配到费用(分)
	AllocPenalty     int        `json:"allocPenalty" binding:""`     // 本次分配到罚息(分)
	Status           int        `json:"status" binding:""`           // 流水状态：1成功 0失败 2冲正/撤销
	Remark           string     `json:"remark" binding:""`           // 备注
}

// UpdateLoanRepaymentTransactionsByIDRequest request params
type UpdateLoanRepaymentTransactionsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id
	// 回款流水ID
	DisbursementID   int64      `json:"disbursementID" binding:""`   // 关联放款单 loan_disbursements.id
	ScheduleID       int64      `json:"scheduleID" binding:""`       // 关联期次 loan_repayment_schedules.id(可空：先入账后分配/未分期)
	CollectChannelID int64      `json:"collectChannelID" binding:""` // 回款渠道(代收) loan_payment_channels.id
	CollectOrderNo   string     `json:"collectOrderNo" binding:""`   // 回款订单号/三方代收单号(商户单号)
	PayRef           string     `json:"payRef" binding:""`           // 支付渠道流水号/交易号(三方transaction id)
	PayAmount        int        `json:"payAmount" binding:""`        // 本次回款金额(分)
	PayMethod        string     `json:"payMethod" binding:""`        // 回款方式(如 BANK_TRANSFER/CARD/WALLET/CASH)
	PaidAt           *time.Time `json:"paidAt" binding:""`           // 回款时间(交易成功时间)
	AllocPrincipal   int        `json:"allocPrincipal" binding:""`   // 本次分配到本金(分)
	AllocInterest    int        `json:"allocInterest" binding:""`    // 本次分配到利息(分)
	AllocFee         int        `json:"allocFee" binding:""`         // 本次分配到费用(分)
	AllocPenalty     int        `json:"allocPenalty" binding:""`     // 本次分配到罚息(分)
	Status           int        `json:"status" binding:""`           // 流水状态：1成功 0失败 2冲正/撤销
	Remark           string     `json:"remark" binding:""`           // 备注
}

// LoanRepaymentTransactionsObjDetail detail
type LoanRepaymentTransactionsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id
	// 回款流水ID
	DisbursementID   int64      `json:"disbursementID"`   // 关联放款单 loan_disbursements.id
	ScheduleID       int64      `json:"scheduleID"`       // 关联期次 loan_repayment_schedules.id(可空：先入账后分配/未分期)
	CollectChannelID int64      `json:"collectChannelID"` // 回款渠道(代收) loan_payment_channels.id
	CollectOrderNo   string     `json:"collectOrderNo"`   // 回款订单号/三方代收单号(商户单号)
	PayRef           string     `json:"payRef"`           // 支付渠道流水号/交易号(三方transaction id)
	PayAmount        int        `json:"payAmount"`        // 本次回款金额(分)
	PayMethod        string     `json:"payMethod"`        // 回款方式(如 BANK_TRANSFER/CARD/WALLET/CASH)
	PaidAt           *time.Time `json:"paidAt"`           // 回款时间(交易成功时间)
	AllocPrincipal   int        `json:"allocPrincipal"`   // 本次分配到本金(分)
	AllocInterest    int        `json:"allocInterest"`    // 本次分配到利息(分)
	AllocFee         int        `json:"allocFee"`         // 本次分配到费用(分)
	AllocPenalty     int        `json:"allocPenalty"`     // 本次分配到罚息(分)
	Status           int        `json:"status"`           // 流水状态：1成功 0失败 2冲正/撤销
	Remark           string     `json:"remark"`           // 备注
	CreatedAt        *time.Time `json:"createdAt"`        // 创建时间
	UpdatedAt        *time.Time `json:"updatedAt"`        // 更新时间
}

// CreateLoanRepaymentTransactionsReply only for api docs
type CreateLoanRepaymentTransactionsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateLoanRepaymentTransactionsByIDReply only for api docs
type UpdateLoanRepaymentTransactionsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetLoanRepaymentTransactionsByIDReply only for api docs
type GetLoanRepaymentTransactionsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanRepaymentTransactions LoanRepaymentTransactionsObjDetail `json:"loanRepaymentTransactions"`
	} `json:"data"` // return data
}

// DeleteLoanRepaymentTransactionsByIDReply only for api docs
type DeleteLoanRepaymentTransactionsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// DeleteLoanRepaymentTransactionssByIDsReply only for api docs
type DeleteLoanRepaymentTransactionssByIDsReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// ListLoanRepaymentTransactionssRequest request params
type ListLoanRepaymentTransactionssRequest struct {
	query.Params
}

// ListLoanRepaymentTransactionssReply only for api docs
type ListLoanRepaymentTransactionssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanRepaymentTransactionss []LoanRepaymentTransactionsObjDetail `json:"loanRepaymentTransactionss"`
	} `json:"data"` // return data
}

// DeleteLoanRepaymentTransactionssByIDsRequest request params
type DeleteLoanRepaymentTransactionssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetLoanRepaymentTransactionsByConditionRequest request params
type GetLoanRepaymentTransactionsByConditionRequest struct {
	query.Conditions
}

// GetLoanRepaymentTransactionsByConditionReply only for api docs
type GetLoanRepaymentTransactionsByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanRepaymentTransactions LoanRepaymentTransactionsObjDetail `json:"loanRepaymentTransactions"`
	} `json:"data"` // return data
}

// ListLoanRepaymentTransactionssByIDsRequest request params
type ListLoanRepaymentTransactionssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListLoanRepaymentTransactionssByIDsReply only for api docs
type ListLoanRepaymentTransactionssByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanRepaymentTransactionss []LoanRepaymentTransactionsObjDetail `json:"loanRepaymentTransactionss"`
	} `json:"data"` // return data
}
