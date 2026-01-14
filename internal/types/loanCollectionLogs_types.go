package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateLoanCollectionLogsRequest request params
type CreateLoanCollectionLogsRequest struct {
	CaseID          int64      `json:"caseID" binding:""`          // 关联催收任务 loan_collection_cases.id
	CollectorUserID int64      `json:"collectorUserID" binding:""` // 催收人员 loan_users.id
	ActionType      string     `json:"actionType" binding:""`      // 动作类型(如 CALL/SMS/VISIT/OTHER，可选)
	Content         string     `json:"content" binding:""`         // 跟进内容/备注(例如用户承诺3天内还款)
	NextFollowUpAt  *time.Time `json:"nextFollowUpAt" binding:""`  // 下次跟进时间(可选)
}

// UpdateLoanCollectionLogsByIDRequest request params
type UpdateLoanCollectionLogsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id
	// 跟进记录ID
	CaseID          int64      `json:"caseID" binding:""`          // 关联催收任务 loan_collection_cases.id
	CollectorUserID int64      `json:"collectorUserID" binding:""` // 催收人员 loan_users.id
	ActionType      string     `json:"actionType" binding:""`      // 动作类型(如 CALL/SMS/VISIT/OTHER，可选)
	Content         string     `json:"content" binding:""`         // 跟进内容/备注(例如用户承诺3天内还款)
	NextFollowUpAt  *time.Time `json:"nextFollowUpAt" binding:""`  // 下次跟进时间(可选)
}

// LoanCollectionLogsObjDetail detail
type LoanCollectionLogsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id
	// 跟进记录ID
	CaseID          int64      `json:"caseID"`          // 关联催收任务 loan_collection_cases.id
	CollectorUserID int64      `json:"collectorUserID"` // 催收人员 loan_users.id
	ActionType      string     `json:"actionType"`      // 动作类型(如 CALL/SMS/VISIT/OTHER，可选)
	Content         string     `json:"content"`         // 跟进内容/备注(例如用户承诺3天内还款)
	NextFollowUpAt  *time.Time `json:"nextFollowUpAt"`  // 下次跟进时间(可选)
	CreatedAt       *time.Time `json:"createdAt"`       // 创建时间
}

// CreateLoanCollectionLogsReply only for api docs
type CreateLoanCollectionLogsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateLoanCollectionLogsByIDReply only for api docs
type UpdateLoanCollectionLogsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetLoanCollectionLogsByIDReply only for api docs
type GetLoanCollectionLogsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanCollectionLogs LoanCollectionLogsObjDetail `json:"loanCollectionLogs"`
	} `json:"data"` // return data
}

// DeleteLoanCollectionLogsByIDReply only for api docs
type DeleteLoanCollectionLogsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// DeleteLoanCollectionLogssByIDsReply only for api docs
type DeleteLoanCollectionLogssByIDsReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// ListLoanCollectionLogssRequest request params
type ListLoanCollectionLogssRequest struct {
	query.Params
}

// ListLoanCollectionLogssReply only for api docs
type ListLoanCollectionLogssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanCollectionLogss []LoanCollectionLogsObjDetail `json:"loanCollectionLogss"`
	} `json:"data"` // return data
}

// DeleteLoanCollectionLogssByIDsRequest request params
type DeleteLoanCollectionLogssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetLoanCollectionLogsByConditionRequest request params
type GetLoanCollectionLogsByConditionRequest struct {
	query.Conditions
}

// GetLoanCollectionLogsByConditionReply only for api docs
type GetLoanCollectionLogsByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanCollectionLogs LoanCollectionLogsObjDetail `json:"loanCollectionLogs"`
	} `json:"data"` // return data
}

// ListLoanCollectionLogssByIDsRequest request params
type ListLoanCollectionLogssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListLoanCollectionLogssByIDsReply only for api docs
type ListLoanCollectionLogssByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanCollectionLogss []LoanCollectionLogsObjDetail `json:"loanCollectionLogss"`
	} `json:"data"` // return data
}
