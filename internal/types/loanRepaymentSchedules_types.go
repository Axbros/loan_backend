package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateLoanRepaymentSchedulesRequest request params
type CreateLoanRepaymentSchedulesRequest struct {
	DisbursementID int64      `json:"disbursementID" binding:""` // 关联放款单 loan_disbursements.id
	InstallmentNo  int        `json:"installmentNo" binding:""`  // 期次(从1开始)
	DueDate        *time.Time `json:"dueDate" binding:""`        // 应还日期
	PrincipalDue   int        `json:"principalDue" binding:""`   // 应还本金(建议统一单位：分)
	InterestDue    int        `json:"interestDue" binding:""`    // 应还利息(分)
	FeeDue         int        `json:"feeDue" binding:""`         // 应还费用(分)
	PenaltyDue     int        `json:"penaltyDue" binding:""`     // 应还罚息(分，逾期产生)
	TotalDue       int        `json:"totalDue" binding:""`       // 本期应还总额=本金+利息+费用+罚息(分)
	PaidPrincipal  int        `json:"paidPrincipal" binding:""`  // 已还本金(分)
	PaidInterest   int        `json:"paidInterest" binding:""`   // 已还利息(分)
	PaidFee        int        `json:"paidFee" binding:""`        // 已还费用(分)
	PaidPenalty    int        `json:"paidPenalty" binding:""`    // 已还罚息(分)
	PaidTotal      int        `json:"paidTotal" binding:""`      // 已还总额(分)
	Status         int        `json:"status" binding:""`         // 期次状态：0未还清 1已还清 2逾期
	LastPaidAt     *time.Time `json:"lastPaidAt" binding:""`     // 最近一次还款时间
	SettledAt      *time.Time `json:"settledAt" binding:""`      // 结清时间(本期还清时)
}

// UpdateLoanRepaymentSchedulesByIDRequest request params
type UpdateLoanRepaymentSchedulesByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id
	// 还款计划ID(期次记录)
	DisbursementID int64      `json:"disbursementID" binding:""` // 关联放款单 loan_disbursements.id
	InstallmentNo  int        `json:"installmentNo" binding:""`  // 期次(从1开始)
	DueDate        *time.Time `json:"dueDate" binding:""`        // 应还日期
	PrincipalDue   int        `json:"principalDue" binding:""`   // 应还本金(建议统一单位：分)
	InterestDue    int        `json:"interestDue" binding:""`    // 应还利息(分)
	FeeDue         int        `json:"feeDue" binding:""`         // 应还费用(分)
	PenaltyDue     int        `json:"penaltyDue" binding:""`     // 应还罚息(分，逾期产生)
	TotalDue       int        `json:"totalDue" binding:""`       // 本期应还总额=本金+利息+费用+罚息(分)
	PaidPrincipal  int        `json:"paidPrincipal" binding:""`  // 已还本金(分)
	PaidInterest   int        `json:"paidInterest" binding:""`   // 已还利息(分)
	PaidFee        int        `json:"paidFee" binding:""`        // 已还费用(分)
	PaidPenalty    int        `json:"paidPenalty" binding:""`    // 已还罚息(分)
	PaidTotal      int        `json:"paidTotal" binding:""`      // 已还总额(分)
	Status         int        `json:"status" binding:""`         // 期次状态：0未还清 1已还清 2逾期
	LastPaidAt     *time.Time `json:"lastPaidAt" binding:""`     // 最近一次还款时间
	SettledAt      *time.Time `json:"settledAt" binding:""`      // 结清时间(本期还清时)
}

// LoanRepaymentSchedulesObjDetail detail
type LoanRepaymentSchedulesObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id
	// 还款计划ID(期次记录)
	DisbursementID int64      `json:"disbursementID"` // 关联放款单 loan_disbursements.id
	InstallmentNo  int        `json:"installmentNo"`  // 期次(从1开始)
	DueDate        *time.Time `json:"dueDate"`        // 应还日期
	PrincipalDue   int        `json:"principalDue"`   // 应还本金(建议统一单位：分)
	InterestDue    int        `json:"interestDue"`    // 应还利息(分)
	FeeDue         int        `json:"feeDue"`         // 应还费用(分)
	PenaltyDue     int        `json:"penaltyDue"`     // 应还罚息(分，逾期产生)
	TotalDue       int        `json:"totalDue"`       // 本期应还总额=本金+利息+费用+罚息(分)
	PaidPrincipal  int        `json:"paidPrincipal"`  // 已还本金(分)
	PaidInterest   int        `json:"paidInterest"`   // 已还利息(分)
	PaidFee        int        `json:"paidFee"`        // 已还费用(分)
	PaidPenalty    int        `json:"paidPenalty"`    // 已还罚息(分)
	PaidTotal      int        `json:"paidTotal"`      // 已还总额(分)
	Status         int        `json:"status"`         // 期次状态：0未还清 1已还清 2逾期
	LastPaidAt     *time.Time `json:"lastPaidAt"`     // 最近一次还款时间
	SettledAt      *time.Time `json:"settledAt"`      // 结清时间(本期还清时)
	CreatedAt      *time.Time `json:"createdAt"`      // 创建时间
	UpdatedAt      *time.Time `json:"updatedAt"`      // 更新时间
}

// CreateLoanRepaymentSchedulesReply only for api docs
type CreateLoanRepaymentSchedulesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateLoanRepaymentSchedulesByIDReply only for api docs
type UpdateLoanRepaymentSchedulesByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetLoanRepaymentSchedulesByIDReply only for api docs
type GetLoanRepaymentSchedulesByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanRepaymentSchedules LoanRepaymentSchedulesObjDetail `json:"loanRepaymentSchedules"`
	} `json:"data"` // return data
}

// DeleteLoanRepaymentSchedulesByIDReply only for api docs
type DeleteLoanRepaymentSchedulesByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// DeleteLoanRepaymentSchedulessByIDsReply only for api docs
type DeleteLoanRepaymentSchedulessByIDsReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// ListLoanRepaymentSchedulessRequest request params
type ListLoanRepaymentSchedulessRequest struct {
	query.Params
}

// ListLoanRepaymentSchedulessReply only for api docs
type ListLoanRepaymentSchedulessReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanRepaymentScheduless []LoanRepaymentSchedulesObjDetail `json:"loanRepaymentScheduless"`
	} `json:"data"` // return data
}

// DeleteLoanRepaymentSchedulessByIDsRequest request params
type DeleteLoanRepaymentSchedulessByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetLoanRepaymentSchedulesByConditionRequest request params
type GetLoanRepaymentSchedulesByConditionRequest struct {
	query.Conditions
}

// GetLoanRepaymentSchedulesByConditionReply only for api docs
type GetLoanRepaymentSchedulesByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanRepaymentSchedules LoanRepaymentSchedulesObjDetail `json:"loanRepaymentSchedules"`
	} `json:"data"` // return data
}

// ListLoanRepaymentSchedulessByIDsRequest request params
type ListLoanRepaymentSchedulessByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListLoanRepaymentSchedulessByIDsReply only for api docs
type ListLoanRepaymentSchedulessByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanRepaymentScheduless []LoanRepaymentSchedulesObjDetail `json:"loanRepaymentScheduless"`
	} `json:"data"` // return data
}
