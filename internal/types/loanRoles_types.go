package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateLoanRolesRequest request params
type CreateLoanRolesRequest struct {
	Code      string `json:"code" binding:""`
	Name      string `json:"name" binding:""`
	DataScope string `json:"dataScope" binding:""`
	Status    int    `json:"status" binding:""`
}

// UpdateLoanRolesByIDRequest request params
type UpdateLoanRolesByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	Code      string `json:"code" binding:""`
	Name      string `json:"name" binding:""`
	DataScope string `json:"dataScope" binding:""`
	Status    int    `json:"status" binding:""`
}

// LoanRolesObjDetail detail
type LoanRolesObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	Code      string     `json:"code"`
	Name      string     `json:"name"`
	DataScope string     `json:"dataScope"`
	Status    int        `json:"status"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

// CreateLoanRolesReply only for api docs
type CreateLoanRolesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateLoanRolesByIDReply only for api docs
type UpdateLoanRolesByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetLoanRolesByIDReply only for api docs
type GetLoanRolesByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanRoles LoanRolesObjDetail `json:"loanRoles"`
	} `json:"data"` // return data
}

// DeleteLoanRolesByIDReply only for api docs
type DeleteLoanRolesByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// DeleteLoanRolessByIDsReply only for api docs
type DeleteLoanRolessByIDsReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// ListLoanRolessRequest request params
type ListLoanRolessRequest struct {
	query.Params
}

// ListLoanRolessReply only for api docs
type ListLoanRolessReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanRoless []LoanRolesObjDetail `json:"loanRoless"`
	} `json:"data"` // return data
}

// DeleteLoanRolessByIDsRequest request params
type DeleteLoanRolessByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetLoanRolesByConditionRequest request params
type GetLoanRolesByConditionRequest struct {
	query.Conditions
}

// GetLoanRolesByConditionReply only for api docs
type GetLoanRolesByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanRoles LoanRolesObjDetail `json:"loanRoles"`
	} `json:"data"` // return data
}

// ListLoanRolessByIDsRequest request params
type ListLoanRolessByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListLoanRolessByIDsReply only for api docs
type ListLoanRolessByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanRoless []LoanRolesObjDetail `json:"loanRoless"`
	} `json:"data"` // return data
}
