package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateLoanLoginAuditRequest request params
type CreateLoanLoginAuditRequest struct {
	UserID    int64  `json:"userID" binding:""`
	LoginType string `json:"loginType" binding:""`
	IP        string `json:"ip" binding:""`
	UserAgent string `json:"userAgent" binding:""`
	Success   int    `json:"success" binding:""`
}

// UpdateLoanLoginAuditByIDRequest request params
type UpdateLoanLoginAuditByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	UserID    int64  `json:"userID" binding:""`
	LoginType string `json:"loginType" binding:""`
	IP        string `json:"ip" binding:""`
	UserAgent string `json:"userAgent" binding:""`
	Success   int    `json:"success" binding:""`
}

// LoanLoginAuditObjDetail detail
type LoanLoginAuditObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	UserID    int64      `json:"userID"`
	LoginType string     `json:"loginType"`
	IP        string     `json:"ip"`
	UserAgent string     `json:"userAgent"`
	Success   int        `json:"success"`
	CreatedAt *time.Time `json:"createdAt"`
}

// CreateLoanLoginAuditReply only for api docs
type CreateLoanLoginAuditReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateLoanLoginAuditByIDReply only for api docs
type UpdateLoanLoginAuditByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetLoanLoginAuditByIDReply only for api docs
type GetLoanLoginAuditByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanLoginAudit LoanLoginAuditObjDetail `json:"loanLoginAudit"`
	} `json:"data"` // return data
}

// DeleteLoanLoginAuditByIDReply only for api docs
type DeleteLoanLoginAuditByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// DeleteLoanLoginAuditsByIDsReply only for api docs
type DeleteLoanLoginAuditsByIDsReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// ListLoanLoginAuditsRequest request params
type ListLoanLoginAuditsRequest struct {
	query.Params
}

// ListLoanLoginAuditsReply only for api docs
type ListLoanLoginAuditsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanLoginAudits []LoanLoginAuditObjDetail `json:"loanLoginAudits"`
	} `json:"data"` // return data
}

// DeleteLoanLoginAuditsByIDsRequest request params
type DeleteLoanLoginAuditsByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetLoanLoginAuditByConditionRequest request params
type GetLoanLoginAuditByConditionRequest struct {
	query.Conditions
}

// GetLoanLoginAuditByConditionReply only for api docs
type GetLoanLoginAuditByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanLoginAudit LoanLoginAuditObjDetail `json:"loanLoginAudit"`
	} `json:"data"` // return data
}

// ListLoanLoginAuditsByIDsRequest request params
type ListLoanLoginAuditsByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListLoanLoginAuditsByIDsReply only for api docs
type ListLoanLoginAuditsByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanLoginAudits []LoanLoginAuditObjDetail `json:"loanLoginAudits"`
	} `json:"data"` // return data
}
