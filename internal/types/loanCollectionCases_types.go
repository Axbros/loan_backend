package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
	"github.com/shopspring/decimal"
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

type CreateLoanCollectionCasesAssignRequest struct {
	CollectorUserID uint64   `json:"collector_user_id" binding:""`
	ScheduleIDs     []uint64 `json:"schedule_ids" binding:""`
}

// types/loan_collection_cases_table.go
type LoanCollectionCasesObjTable struct {
	ID         uint64 `json:"id" gorm:"column:id"`
	ScheduleID uint64 `json:"schedule_id" gorm:"column:schedule_id"`
	BaseinfoID uint64 `json:"baseinfo_id" gorm:"column:baseinfo_id"`

	FirstName  string `json:"first_name" gorm:"column:first_name"`
	SecondName string `json:"second_name" gorm:"column:second_name"`
	Age        int    `json:"age" gorm:"column:age"`
	Gender     string `json:"gender" gorm:"column:gender"`
	IDType     string `json:"id_type" gorm:"column:id_type"`
	IDNumber   string `json:"id_number" gorm:"column:id_number"`
	Mobile     string `json:"mobile" gorm:"column:mobile"`

	Priority      int        `json:"priority" gorm:"column:priority"`
	Status        int        `json:"status" gorm:"column:status"`
	CompletedAt   *time.Time `json:"completed_at" gorm:"column:completed_at"`
	CompletedNote string     `json:"completed_note" gorm:"column:completed_note"`

	DueDate   time.Time       `json:"due_date" gorm:"column:due_date"`
	NetAmount decimal.Decimal `json:"net_amount" gorm:"column:net_amount"` // 或 float64
	TotalDue  decimal.Decimal `json:"total_due" gorm:"column:total_due"`   // 或 float64
	PaidTotal decimal.Decimal `json:"paid_total" gorm:"column:paid_total"` // 或 float64

	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at"`
	CollectorName  string    `json:"collector_name" gorm:"column:collector_name"`
	AssignedByName *string   `json:"assigned_by_name" gorm:"column:assigned_by_name"`
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
