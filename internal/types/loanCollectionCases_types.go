package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateLoanCollectionCasesRequest request params
type CreateLoanCollectionCasesRequest struct {
	DisbursementID   int64      `json:"disbursementID" binding:""`   // 关联放款单 loan_disbursements.id
	ScheduleID       int64      `json:"scheduleID" binding:""`       // 关联逾期期次 loan_repayment_schedules.id(按期催收可用，整单催收可为空)
	CollectorUserID  int64      `json:"collectorUserID" binding:""`  // 催收人员 loan_users.id
	AssignedByUserID int64      `json:"assignedByUserID" binding:""` // 分配人(管理员) loan_users.id
	AssignedAt       *time.Time `json:"assignedAt" binding:""`       // 分配时间
	Priority         int        `json:"priority" binding:""`         // 优先级：1高 2中 3低
	Status           int        `json:"status" binding:""`           // 任务状态：0待处理 1跟进中 2已完成 3已取消
	DueAmount        int        `json:"dueAmount" binding:""`        // 逾期应还金额快照(分，可选，用于列表展示)
	OverdueDays      int        `json:"overdueDays" binding:""`      // 逾期天数快照(可选，用于列表展示)
	CompletedAt      *time.Time `json:"completedAt" binding:""`      // 完成时间(点击完成时)
	CompletedNote    string     `json:"completedNote" binding:""`    // 完成备注(例如用户承诺X天内还款)
}

// UpdateLoanCollectionCasesByIDRequest request params
type UpdateLoanCollectionCasesByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id
	// 催收任务ID
	DisbursementID   int64      `json:"disbursementID" binding:""`   // 关联放款单 loan_disbursements.id
	ScheduleID       int64      `json:"scheduleID" binding:""`       // 关联逾期期次 loan_repayment_schedules.id(按期催收可用，整单催收可为空)
	CollectorUserID  int64      `json:"collectorUserID" binding:""`  // 催收人员 loan_users.id
	AssignedByUserID int64      `json:"assignedByUserID" binding:""` // 分配人(管理员) loan_users.id
	AssignedAt       *time.Time `json:"assignedAt" binding:""`       // 分配时间
	Priority         int        `json:"priority" binding:""`         // 优先级：1高 2中 3低
	Status           int        `json:"status" binding:""`           // 任务状态：0待处理 1跟进中 2已完成 3已取消
	DueAmount        int        `json:"dueAmount" binding:""`        // 逾期应还金额快照(分，可选，用于列表展示)
	OverdueDays      int        `json:"overdueDays" binding:""`      // 逾期天数快照(可选，用于列表展示)
	CompletedAt      *time.Time `json:"completedAt" binding:""`      // 完成时间(点击完成时)
	CompletedNote    string     `json:"completedNote" binding:""`    // 完成备注(例如用户承诺X天内还款)
}

// LoanCollectionCasesObjDetail detail
type LoanCollectionCasesObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id
	// 催收任务ID
	DisbursementID   int64      `json:"disbursementID"`   // 关联放款单 loan_disbursements.id
	ScheduleID       int64      `json:"scheduleID"`       // 关联逾期期次 loan_repayment_schedules.id(按期催收可用，整单催收可为空)
	CollectorUserID  int64      `json:"collectorUserID"`  // 催收人员 loan_users.id
	AssignedByUserID int64      `json:"assignedByUserID"` // 分配人(管理员) loan_users.id
	AssignedAt       *time.Time `json:"assignedAt"`       // 分配时间
	Priority         int        `json:"priority"`         // 优先级：1高 2中 3低
	Status           int        `json:"status"`           // 任务状态：0待处理 1跟进中 2已完成 3已取消
	DueAmount        int        `json:"dueAmount"`        // 逾期应还金额快照(分，可选，用于列表展示)
	OverdueDays      int        `json:"overdueDays"`      // 逾期天数快照(可选，用于列表展示)
	CompletedAt      *time.Time `json:"completedAt"`      // 完成时间(点击完成时)
	CompletedNote    string     `json:"completedNote"`    // 完成备注(例如用户承诺X天内还款)
	CreatedAt        *time.Time `json:"createdAt"`        // 创建时间
	UpdatedAt        *time.Time `json:"updatedAt"`        // 更新时间
}

// CreateLoanCollectionCasesReply only for api docs
type CreateLoanCollectionCasesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateLoanCollectionCasesByIDReply only for api docs
type UpdateLoanCollectionCasesByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetLoanCollectionCasesByIDReply only for api docs
type GetLoanCollectionCasesByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanCollectionCases LoanCollectionCasesObjDetail `json:"loanCollectionCases"`
	} `json:"data"` // return data
}

// DeleteLoanCollectionCasesByIDReply only for api docs
type DeleteLoanCollectionCasesByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// DeleteLoanCollectionCasessByIDsReply only for api docs
type DeleteLoanCollectionCasessByIDsReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// ListLoanCollectionCasessRequest request params
type ListLoanCollectionCasessRequest struct {
	query.Params
}

// ListLoanCollectionCasessReply only for api docs
type ListLoanCollectionCasessReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanCollectionCasess []LoanCollectionCasesObjDetail `json:"loanCollectionCasess"`
	} `json:"data"` // return data
}

// DeleteLoanCollectionCasessByIDsRequest request params
type DeleteLoanCollectionCasessByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetLoanCollectionCasesByConditionRequest request params
type GetLoanCollectionCasesByConditionRequest struct {
	query.Conditions
}

// GetLoanCollectionCasesByConditionReply only for api docs
type GetLoanCollectionCasesByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanCollectionCases LoanCollectionCasesObjDetail `json:"loanCollectionCases"`
	} `json:"data"` // return data
}

// ListLoanCollectionCasessByIDsRequest request params
type ListLoanCollectionCasessByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListLoanCollectionCasessByIDsReply only for api docs
type ListLoanCollectionCasessByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanCollectionCasess []LoanCollectionCasesObjDetail `json:"loanCollectionCasess"`
	} `json:"data"` // return data
}
