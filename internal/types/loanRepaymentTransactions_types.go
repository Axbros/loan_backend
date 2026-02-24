package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateLoanRepaymentTransactionsRequest request params
type CreateLoanRepaymentTransactionsRequest struct {
	ScheduleID       uint64 `json:"scheduleID" binding:""` // 关联期次 loan_repayment_schedules.id(可空：先入账后分配/未分期)
	CollectChannelID int64  `json:"collectChannelID" binding:""`
	PayAmount        int    `json:"payAmount" binding:""`      // 本次回款金额(分)
	PayMethod        string `json:"payMethod" binding:""`      // 回款方式(如 BANK_TRANSFER/WALLET)
	AllocPrincipal   int    `json:"allocPrincipal" binding:""` // 本次分配到本金(分)
	AllocInterest    int    `json:"allocInterest" binding:""`  // 本次分配到利息(分)
	AllocFee         int    `json:"allocFee" binding:""`       // 本次分配到费用(分)
	AllocPenalty     int    `json:"allocPenalty" binding:""`   // 本次分配到罚息(分)
	VoucherFileName  string `json:"voucherFileName" binding:""`
	MfaCode          string `json:"mfaCode" binding:""`
	Remark           string `json:"remark" binding:""` // 备注
}

type DetailByScheduleIDRequest struct {
	ScheduleID uint64 `json:"schedule_id" binding:""`
}

// RepaymentScheduleDetail 还款计划详情查询结果结构体
// 对应你的多表关联查询结果
type RepaymentScheduleDetail struct {
	FirstName         string     `gorm:"column:first_name" json:"firstName"`                 // 借款人名字
	SecondName        string     `gorm:"column:second_name" json:"secondName"`               // 借款人姓氏
	Age               int        `gorm:"column:age" json:"age"`                              // 借款人年龄
	Gender            string     `gorm:"column:gender" json:"gender"`                        // 借款人性别
	IDType            string     `gorm:"column:id_type" json:"idType"`                       // 证件类型
	IDNumber          string     `gorm:"column:id_number" json:"idNumber"`                   // 证件号码
	ApplicationAmount int64      `gorm:"column:application_amount" json:"applicationAmount"` // 申请金额（分）
	NetAmount         int64      `gorm:"column:net_amount" json:"netAmount"`                 // 放款净金额（分）
	PayoutOrderNo     string     `gorm:"column:payout_order_no" json:"payoutOrderNo"`        // 放款订单号
	DisbursedAt       *time.Time `gorm:"column:disbursed_at" json:"disbursedAt"`             // 放款时间
	DueDate           *time.Time `gorm:"column:due_date" json:"dueDate"`                     // 应还日期
	PaidTotal         int64      `gorm:"column:paid_total" json:"paidTotal"`                 // 已还总额（分）
	TotalDue          int64      `gorm:"column:total_due" json:"totalDue"`                   // 应还总额（分）
	ChannelName       string     `gorm:"column:name" json:"channelName"`                     // 支付渠道名称
	PayoutFeeRate     int64      `gorm:"column:payout_fee_rate" json:"payoutFeeRate"`        // 手续费率
}

// UpdateLoanRepaymentTransactionsByIDRequest request params
type UpdateLoanRepaymentTransactionsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id
	// 回款流水ID
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

type LoanRepaymentTransactionsHistory struct {
	ID                 uint64     `json:"id"`                 // 流水ID（对应t.id）
	CollectOrderNo     string     `json:"collectOrderNo"`     // 回款订单号（对应t.collect_order_no）
	PayAmount          int        `json:"payAmount"`          // 回款金额（对应t.pay_amount）
	AllocPrincipal     int        `json:"allocPrincipal"`     // 分配本金（对应t.alloc_principal）
	AllocInterest      int        `json:"allocInterest"`      // 分配利息（对应t.alloc_interest）
	AllocFee           int        `json:"allocFee"`           // 分配费用（对应t.alloc_fee）
	AllocPenalty       int        `json:"allocPenalty"`       // 分配罚息（对应t.alloc_penalty）
	Status             int        `json:"status"`             // 状态（对应t.status）
	Remark             string     `json:"remark"`             // 备注（对应t.remark）
	VoucherFileName    string     `json:"voucherFileName"`    // 凭证文件名（对应t.voucher_file_name）
	CreatedAt          *time.Time `json:"createdAt"`          // 创建时间（对应t.created_at）
	CollectChannelName string     `json:"collectChannelName"` // 回款渠道名称（对应c.name）
	CreatedBy          string     `json:"createdBy"`          // 创建人用户名（对应u.username）
	// 保留原有关联字段（如果需要）
	ScheduleID       int64 `json:"scheduleID"`       // 关联期次ID
	CollectChannelID int64 `json:"collectChannelID"` // 回款渠道ID
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
