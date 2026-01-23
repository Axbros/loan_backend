package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateLoanAuditsRequest request params
type CreateLoanAuditsRequest struct {
	BaseinfoID    int    `json:"baseinfoID" binding:""`    // 关联申请单 loan_baseinfo.id
	AuditResult   int    `json:"auditResult" binding:""`   // 审核结果：1通过 -1拒绝
	AuditComment  string `json:"auditComment" binding:""`  // 审核备注/原因
	AuditorUserID int64  `json:"auditorUserID" binding:""` // 审核人员(loan_users.id)
}

// UpdateLoanAuditsByIDRequest request params
type UpdateLoanAuditsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id
	// 审核记录ID
	BaseinfoID    int    `json:"baseinfoID" binding:""`    // 关联申请单 loan_baseinfo.id
	AuditResult   int    `json:"auditResult" binding:""`   // 审核结果：1通过 -1拒绝
	AuditComment  string `json:"auditComment" binding:""`  // 审核备注/原因
	AuditorUserID int64  `json:"auditorUserID" binding:""` // 审核人员(loan_users.id)
}

// LoanAuditsObjDetail detail
type LoanAuditsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id
	// 审核记录ID
	BaseinfoID    int    `json:"baseinfoID"`    // 关联申请单 loan_baseinfo.id
	AuditResult   int    `json:"auditResult"`   // 审核结果：1通过 -1拒绝
	AuditComment  string `json:"auditComment"`  // 审核备注/原因
	AuditorUserID int64  `json:"auditorUserID"` // 审核人员(loan_users.id)
	// 新增字段：审核人员真实姓名（非数据库字段，仅用于返回）
	AuditorName string     `json:"auditorName"` // gorm:"-" 表示不映射数据库字段
	CreatedAt   *time.Time `json:"createdAt"`   // 审核时间(即审核通过/拒绝时间)
}

type LoanAuditDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id
	// 审核记录ID
	BaseinfoID   int    `json:"baseinfoID"`   // 关联申请单 loan_baseinfo.id
	AuditResult  int    `json:"auditResult"`  // 审核结果：1通过 -1拒绝
	AuditComment string `json:"auditComment"` // 审核备注/原因
	//AuditorUserID int64      `json:"auditorUserID"` // 审核人员(loan_users.id)
	AuditType     int        `json:"auditType"`
	AuditUsername string     `gorm:"column:auditor_username" json:"auditUsername"`
	CreatedAt     *time.Time `json:"createdAt"` // 审核时间(即审核通过/拒绝时间)
}

// CreateLoanAuditsReply only for api docs
type CreateLoanAuditsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateLoanAuditsByIDReply only for api docs
type UpdateLoanAuditsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetLoanAuditsByIDReply only for api docs
type GetLoanAuditsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanAudits LoanAuditsObjDetail `json:"loanAudits"`
	} `json:"data"` // return data
}

// DeleteLoanAuditsByIDReply only for api docs
type DeleteLoanAuditsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// DeleteLoanAuditssByIDsReply only for api docs
type DeleteLoanAuditssByIDsReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

type RequestLoanAuditsDetail struct {
	BaseinfoID uint64 `json:"baseinfoID" binding:""`
}

// ListLoanAuditssRequest request params
type ListLoanAuditssRequest struct {
	query.Params
}

// ListLoanAuditssReply only for api docs
type ListLoanAuditssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanAuditss []LoanAuditsObjDetail `json:"loanAuditss"`
	} `json:"data"` // return data
}

// DeleteLoanAuditssByIDsRequest request params
type DeleteLoanAuditssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetLoanAuditsByConditionRequest request params
type GetLoanAuditsByConditionRequest struct {
	query.Conditions
}

// GetLoanAuditsByConditionReply only for api docs
type GetLoanAuditsByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanAudits LoanAuditsObjDetail `json:"loanAudits"`
	} `json:"data"` // return data
}

// ListLoanAuditssByIDsRequest request params
type ListLoanAuditssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListLoanAuditssByIDsReply only for api docs
type ListLoanAuditssByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanAuditss []LoanAuditsObjDetail `json:"loanAuditss"`
	} `json:"data"` // return data
}
