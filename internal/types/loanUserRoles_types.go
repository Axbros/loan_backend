package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateLoanUserRolesRequest request params
type CreateLoanUserRolesRequest struct {
	UserID int64 `json:"userID" binding:""`
	RoleID int64 `json:"roleID" binding:""`
}

// UpdateLoanUserRolesByIDRequest request params
type UpdateLoanUserRolesByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id
	// surrogate id
}

// LoanUserRolesObjDetail detail
type LoanUserRolesObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id
	// surrogate id
	UpdatedAt *time.Time `json:"updatedAt"` // 更新时间
}

// CreateLoanUserRolesReply only for api docs
type CreateLoanUserRolesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateLoanUserRolesByIDReply only for api docs
type UpdateLoanUserRolesByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetLoanUserRolesByIDReply only for api docs
type GetLoanUserRolesByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanUserRoles LoanUserRolesObjDetail `json:"loanUserRoles"`
	} `json:"data"` // return data
}

// DeleteLoanUserRolesByIDReply only for api docs
type DeleteLoanUserRolesByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// DeleteLoanUserRolessByIDsReply only for api docs
type DeleteLoanUserRolessByIDsReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// ListLoanUserRolessRequest request params
type ListLoanUserRolessRequest struct {
	query.Params
}

// ListLoanUserRolessReply only for api docs
type ListLoanUserRolessReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanUserRoless []LoanUserRolesObjDetail `json:"loanUserRoless"`
	} `json:"data"` // return data
}

// DeleteLoanUserRolessByIDsRequest request params
type DeleteLoanUserRolessByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetLoanUserRolesByConditionRequest request params
type GetLoanUserRolesByConditionRequest struct {
	query.Conditions
}

// GetLoanUserRolesByConditionReply only for api docs
type GetLoanUserRolesByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanUserRoles LoanUserRolesObjDetail `json:"loanUserRoles"`
	} `json:"data"` // return data
}

// ListLoanUserRolessByIDsRequest request params
type ListLoanUserRolessByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListLoanUserRolessByIDsReply only for api docs
type ListLoanUserRolessByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanUserRoless []LoanUserRolesObjDetail `json:"loanUserRoless"`
	} `json:"data"` // return data
}
