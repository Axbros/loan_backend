package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateLoanDisbursementsRequest request params
type CreateLoanDisbursementsRequest struct {
	BaseinfoID           int        `json:"baseinfoID" binding:""`           // 关联申请单 loan_baseinfo.id
	DisburseAmount       int        `json:"disburseAmount" binding:""`       // 放款金额(单位按你的系统：元/分，建议统一)
	NetAmount            int        `json:"netAmount" binding:""`            // 到账金额(扣除费用后实际到账)
	Status               int        `json:"status" binding:""`               // 放款状态：0待放款 1已放款
	SourceReferrerUserID int64      `json:"sourceReferrerUserID" binding:""` // 用户来源(分享人 loan_users.id，冗余快照，便于查询)
	AuditorUserID        int64      `json:"auditorUserID" binding:""`        // 审核人员(loan_users.id)
	AuditedAt            *time.Time `json:"auditedAt" binding:""`            // 审核通过时间
	PayoutChannelID      int64      `json:"payoutChannelID" binding:""`      // 放款渠道(代付) loan_payment_channels.id
	PayoutOrderNo        string     `json:"payoutOrderNo" binding:""`        // 放款订单号/三方代付单号
	DisbursedAt          *time.Time `json:"disbursedAt" binding:""`          // 放款时间
}

// UpdateLoanDisbursementsByIDRequest request params
type UpdateLoanDisbursementsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id
	// 主键(放款单ID)
	BaseinfoID           int        `json:"baseinfoID" binding:""`           // 关联申请单 loan_baseinfo.id
	DisburseAmount       int        `json:"disburseAmount" binding:""`       // 放款金额(单位按你的系统：元/分，建议统一)
	NetAmount            int        `json:"netAmount" binding:""`            // 到账金额(扣除费用后实际到账)
	Status               int        `json:"status" binding:""`               // 放款状态：0待放款 1已放款
	SourceReferrerUserID int64      `json:"sourceReferrerUserID" binding:""` // 用户来源(分享人 loan_users.id，冗余快照，便于查询)
	AuditorUserID        int64      `json:"auditorUserID" binding:""`        // 审核人员(loan_users.id)
	AuditedAt            *time.Time `json:"auditedAt" binding:""`            // 审核通过时间
	PayoutChannelID      int64      `json:"payoutChannelID" binding:""`      // 放款渠道(代付) loan_payment_channels.id
	PayoutOrderNo        string     `json:"payoutOrderNo" binding:""`        // 放款订单号/三方代付单号
	DisbursedAt          *time.Time `json:"disbursedAt" binding:""`          // 放款时间
}

// LoanDisbursementsObjDetail detail
type LoanDisbursementsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id
	// 主键(放款单ID)
	BaseinfoID           int        `json:"baseinfoID"`           // 关联申请单 loan_baseinfo.id
	DisburseAmount       int        `json:"disburseAmount"`       // 放款金额(单位按你的系统：元/分，建议统一)
	NetAmount            int        `json:"netAmount"`            // 到账金额(扣除费用后实际到账)
	Status               int        `json:"status"`               // 放款状态：0待放款 1已放款
	SourceReferrerUserID int64      `json:"sourceReferrerUserID"` // 用户来源(分享人 loan_users.id，冗余快照，便于查询)
	AuditorUserID        int64      `json:"auditorUserID"`        // 审核人员(loan_users.id)
	AuditedAt            *time.Time `json:"auditedAt"`            // 审核通过时间
	PayoutChannelID      int64      `json:"payoutChannelID"`      // 放款渠道(代付) loan_payment_channels.id
	PayoutOrderNo        string     `json:"payoutOrderNo"`        // 放款订单号/三方代付单号
	DisbursedAt          *time.Time `json:"disbursedAt"`          // 放款时间
	CreatedAt            *time.Time `json:"createdAt"`            // 创建时间(进入待放款时刻)
	UpdatedAt            *time.Time `json:"updatedAt"`            // 更新时间
}

// CreateLoanDisbursementsReply only for api docs
type CreateLoanDisbursementsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

type LoanDisbursementDetailReply struct {
	Disbursement LoanDisbursementsObjDetail `json:"disbursement"`
	Baseinfo     LoanBaseinfoObjDetail      `json:"baseinfo"`
	Audits       LoanAuditsObjDetail        `json:"audits"`
}

// UpdateLoanDisbursementsByIDReply only for api docs
type UpdateLoanDisbursementsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetLoanDisbursementsByIDReply only for api docs
type GetLoanDisbursementsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanDisbursements LoanDisbursementsObjDetail `json:"loanDisbursements"`
	} `json:"data"` // return data
}

// DeleteLoanDisbursementsByIDReply only for api docs
type DeleteLoanDisbursementsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// DeleteLoanDisbursementssByIDsReply only for api docs
type DeleteLoanDisbursementssByIDsReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// ListLoanDisbursementssRequest request params
type ListLoanDisbursementssRequest struct {
	query.Params
}

// ListLoanDisbursementssReply only for api docs
type ListLoanDisbursementssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanDisbursementss []LoanDisbursementsObjDetail `json:"loanDisbursementss"`
	} `json:"data"` // return data
}

// DeleteLoanDisbursementssByIDsRequest request params
type DeleteLoanDisbursementssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetLoanDisbursementsByConditionRequest request params
type GetLoanDisbursementsByConditionRequest struct {
	query.Conditions
}

// GetLoanDisbursementsByConditionReply only for api docs
type GetLoanDisbursementsByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanDisbursements LoanDisbursementsObjDetail `json:"loanDisbursements"`
	} `json:"data"` // return data
}

// ListLoanDisbursementssByIDsRequest request params
type ListLoanDisbursementssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListLoanDisbursementssByIDsReply only for api docs
type ListLoanDisbursementssByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanDisbursementss []LoanDisbursementsObjDetail `json:"loanDisbursementss"`
	} `json:"data"` // return data
}
